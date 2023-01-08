package group

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
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
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	res, err := insert.Exec(name)
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// グループIDを取得
	groupId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	//グループを作ったユーザを管理者にする
	tokenCookie, _ := c.Cookie("token")
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	update, err := db.Prepare("UPDATE users SET group_id = ?, is_manager = true WHERE user_id = ?;")
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// SQLの実行
	_, err = update.Exec(groupId, userId)
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
 }

func Participate(c echo.Context) error {
	groupId := c.FormValue("group_id")
	
	tokenCookie, _ := c.Cookie("token")
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// ユーザをグループに入会させる処理
	update, err := db.Prepare("UPDATE users SET group_id = ? WHERE user_id = ?;")
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer update.Close()

	_, err = update.Exec(groupId, userId)
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// func Select(c echo.Context) error {}

// func Quit(c echo.Context) error {}

// func HostUpdate(c echo.Context) error {}