package group

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
)

 func Register(c echo.Context) error {
	// 値を変数に格納
	name := c.FormValue("group_name")

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// SQL文を定義
	insert, err := db.Prepare("INSERT INTO groups(group_name) VALUES(?);")
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusOK)
	}
	defer insert.Close()

	// SQLの実行
	_, err = insert.Exec(name)
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusOK)
 }

// func Participate(c echo.Context) error {}

// func Select(c echo.Context) error {}

// func Quit(c echo.Context) error {}

// func HostUpdate(c echo.Context) error {}