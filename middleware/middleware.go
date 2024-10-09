package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

var SECRETKEY string

func SetSecretKey(secret_key string) {
	SECRETKEY = secret_key
}

func Sign(data, secret_key []byte) (string, error) {
	mac := hmac.New(sha256.New, secret_key)

	_, err := mac.Write(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func Authenticate(secret_key, data, hash []byte) (bool, error) {
	mac := hmac.New(sha256.New, secret_key)

	_, err := mac.Write(data)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(mac.Sum(nil)) != hex.EncodeToString(hash), nil
}

func CreatePasswordHash(password, salt string) (string, error) {
	passwordBytes := append([]byte(password), []byte(salt)...)

	sha512hasher := sha512.New()
	sha512hasher.Write(passwordBytes)

	return hex.EncodeToString(sha512hasher.Sum(nil)), nil
}
