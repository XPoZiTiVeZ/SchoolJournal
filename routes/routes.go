package routes

import (
	"database/sql"
	"fmt"
	"journal/middleware"
	"journal/models"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

var urlMap map[string]string
var db *sql.DB

func IndexHandlerGET(w http.ResponseWriter, r *http.Request) {
	tpl, _ := gonja.FromFile("./templates/index.html")
	ctx := exec.NewContext(map[string]interface{}{
		"title": "Главная",
		"url":   urlMap,
	})
	tpl.Execute(w, ctx)
}

func LoginHandlerGET(w http.ResponseWriter, r *http.Request) {
	tpl, _ := gonja.FromFile("./templates/account/login.html")
	ctx := exec.NewContext(map[string]interface{}{
		"title": "Вход в аккаунт",
		"url":   urlMap,
	})
	tpl.Execute(w, ctx)
}

func LoginHandlerPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	_, err = mail.ParseAddress(email)
	if err != nil {
		fmt.Println(err)
	}
	user := models.GetUser(email)

	passwordHash, _ := middleware.CreatePasswordHash(password, user.Salt)

	if user.PasswordHash == passwordHash {
		sessionUUID := uuid.NewString()
		signature, _ := middleware.Sign([]byte(sessionUUID), []byte(middleware.SECRETKEY))

		cookie1 := http.Cookie{
			Name:     "sessionToken",
			Value:    signature,
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		cookie2 := http.Cookie{
			Name:     "sessionUUID",
			Value:    sessionUUID,
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie1)
		http.SetCookie(w, &cookie2)

	}

	LoginHandlerGET(w, r)
	// tpl, _ := gonja.FromFile("./templates/account/login.html")
	// ctx := exec.NewContext(map[string]interface{}{
	// 	"title": "Вход в аккаунт",
	// 	"url":   urlMap,
	// })
	// tpl.Execute(w, ctx)
}

func StartServer(url_map map[string]string, db_ *sql.DB) {
	urlMap = url_map
	db = db_

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))
	http.HandleFunc(fmt.Sprint("GET ", urlMap["IndexPage"]), IndexHandlerGET)
	http.HandleFunc(fmt.Sprint("GET ", urlMap["LoginPage"]), LoginHandlerGET)
	http.HandleFunc(fmt.Sprint("POST ", urlMap["LoginPage"]), LoginHandlerPOST)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
