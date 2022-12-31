package user

import (
	"net/http"
	"log"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/hash"
	"src/lib/jwt"
	"src/lib/token"
)

func Register(c echo.Context) error {
	// form-dataの値を変数に格納
	name 	 := c.FormValue("user_name")
	password := c.FormValue("user_password")
	email 	 := c.FormValue("email")

	// パスワードをハッシュ化
	hashedPassword := hash.HashPassword(password)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// SQL文を定義
	insert, err := db.Prepare("INSERT INTO users(user_name, user_password, email) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusOK)
	}
	defer insert.Close()

	// SQLの実行
	res, err := insert.Exec(name, hashedPassword, email)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusOK)
	}

	// ユーザIDを取得
	userId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusOK)
	}

	// トークンを発行
	tokenText := jwt.GetTokenText(userId)
	tokenCookie := token.GetToken(tokenText)
	c.SetCookie(tokenCookie)

	return c.NoContent(http.StatusOK)
}

// func Authentication(c echo.Context) error {}

// func CheckLogin(c echo.Context) error {}