package token

import (
	"net/http"
	"time"
	"fmt"
)

func GetToken(tokenText string) *http.Cookie {
	fmt.Println(tokenText)
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenText
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	return cookie
}