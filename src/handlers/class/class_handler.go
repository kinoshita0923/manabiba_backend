package class

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
)

func Register(c echo.Context) error {
	className	:= c.FormValue("class_name")
	grade 		:= c.FormValue("grade")

	// トークンを変数に格納
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusOK, "No token")
	}
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	var (
		groupId int
	)

	rows, err := db.Query(
		"SELECT group_id FROM users WHERE user_id = ?;",
		userId,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&groupId)

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	// SQL文を定義
	insert, err := db.Prepare("INSERT INTO classes(class_name, grade, group_id) VALUES(?, ?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	_, err = insert.Exec(className, grade, groupId)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}