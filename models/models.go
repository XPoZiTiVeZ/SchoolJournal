package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	Id           string `json:"id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	PasswordHash string
	Salt         string
}

func CreateUserDB() {
	stmt, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users (
			Id 		     TEXT PRIMARY KEY, 
			FirstName	 TEXT,
			LastName	 TEXT,
			Email        TEXT, 
			PasswordHash TEXT, 
			Salt 		 TEXT
		);`,
	)
	stmt.Exec()
	stmt.Close()
}

func GetUser(email string) User {
	user := User{}
	stmt, _ := db.Prepare(`SELECT users.Id, users.FirstName, users.LastName, users.Email, users.PasswordHash, users.Salt FROM users WHERE users.Email = ?;`)
	stmt.QueryRow(email).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
	)

	fmt.Println(user)
	return user
}

type Session struct {
	UserId       string `json:"user_id"`
	SessionToken string `json:"session_token"`
}

func CreateSessionDB() {
	stmt, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS sessions (
			UserId 		 TEXT PRIMARY KEY,
			SessionToken TEXT
		);`,
	)
	stmt.Exec()
}

func AddSession(email, sessionToken string) {

}

func CreateAllDB(db_file string) *sql.DB {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		os.Create(db_file)
		db, _ = sql.Open("sqlite3", db_file)
	}

	CreateUserDB()
	CreateSessionDB()

	return db
}
