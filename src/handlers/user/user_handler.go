package user

import (
	"log"
	"net/http"

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
	insert, err := db.Prepare("INSERT INTO users(user_name, user_password, email) VALUES(?, ?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(500)
	}
	defer insert.Close()

	// SQLの実行
	res, err := insert.Exec(name, hashedPassword, email)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(500)
	}

	// ユーザIDを取得
	userId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return c.NoContent(500)
	}

	// トークンを発行
	tokenText := jwt.GetTokenText(userId)
	tokenCookie := token.GetToken(tokenText)
	c.SetCookie(tokenCookie)

	return c.NoContent(http.StatusOK)
}

func Authentication(c echo.Context) error {
	// form-dataの値を変数に格納	
	email	 := c.FormValue("email")
	password := c.FormValue("password")

	// パスワードをハッシュ化
	hashedPassword := hash.HashPassword(password)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// DBにクエリを送信
	rows, err := db.Query(
		"SELECT user_id, user_password FROM users WHERE email = ?;",
		email,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(500)
	}
	defer rows.Close()

	// 結果を代入する変数を定義
	var (
		user_id       int
		user_password string
	)

	// 結果を代入
	for rows.Next() {
		err := rows.Scan(&user_id, &user_password)

		if err != nil {
			return c.NoContent(500)
		}
	}

	// トークンを発行
	if user_password == hashedPassword {
		tokenText := jwt.GetTokenText(int64(user_id))
		tokenCookie := token.GetToken(tokenText)
		c.SetCookie(tokenCookie)
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(403)
	}
}

func CheckLogin(c echo.Context) error {
	// トークンを変数に格納
	tokenCookie, _ := c.Cookie("token")

	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// DBにクエリを送信
	rows, err := db.Query(
		"SELECT EXISTS(SELECT * FROM users WHERE user_id = ?) AS exist_check;",
		userId,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(500)
	}
	defer rows.Close()

	// 結果を代入する変数を定義
	var (
		exist_check int
	)

	// 結果を代入
	for rows.Next() {
		err := rows.Scan(&exist_check)

		if err != nil {
			return c.NoContent(500)
		}
	}

	// トークンを発行
	if exist_check >= 1 {
		tokenText := jwt.GetTokenText(int64(userId))
		token := token.GetToken(tokenText)
		c.SetCookie(token)
		return c.NoContent(http.StatusOK)
	} else {
		cookie := token.DeleteToken()
		c.SetCookie(cookie)
		return c.NoContent(403)
	}
}