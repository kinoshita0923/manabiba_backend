package token

import (
	"net/http"
	"time"
)

func GenerateCookie(tokenText string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenText
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	return cookie
}

func DeleteCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Expires = time.Now()
	return cookie
}