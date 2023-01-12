package class

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
)

func Register(c echo.Context) error {
	className	:= c.FormValue("class_name")
	grade 		:= c.FormValue("grade")

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// SQL文を定義
	insert, err := db.Prepare("INSERT INTO classes(class_name, grade) VALUES(?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	_, err = insert.Exec(className, grade)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}