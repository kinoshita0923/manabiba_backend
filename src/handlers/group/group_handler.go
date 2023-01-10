package group

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
)

type Group struct {
	GroupId 	int
	GroupName 	string
}

type Groups struct {
	Groups []Group
}

 func Register(c echo.Context) error {
	// 値を変数に格納
	groupName := c.FormValue("group_name")

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
	res, err := insert.Exec(groupName)
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
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusUnauthorized, "No token")
	}
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
	
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusUnauthorized, "No token")
	}
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

func Select(c echo.Context) error {
	groupName := c.QueryParam("group_name")
	fmt.Println(groupName)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	rows, err := db.Query(
		"SELECT group_id, group_name FROM groups WHERE group_name LIKE ?;",
		"%" + groupName + "%",
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rows.Close()

	fmt.Println("test1 success")

	groups := &Groups{}
	for rows.Next() {
		fmt.Println("test2 success")
		group := &Group{}
		err := rows.Scan(&group.GroupId, &group.GroupName)
		fmt.Println(*group)

		if len(groups.Groups) == 0 {
			groups.Groups = []Group{*group}
		} else {
			groups.Groups = append(groups.Groups, *group)
		}

		fmt.Println(*groups)

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	fmt.Println(*groups)

	return c.JSON(http.StatusOK, groups)
}

// func Quit(c echo.Context) error {}

func HostUpdate(c echo.Context) error {
	// 引継ぎ先のユーザのメールアドレスを取得
	email := c.FormValue("email")

	// トークンからユーザIDを取得
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusUnauthorized, "No token")
	}
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	db := database.Connect()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// 結果を代入する変数を宣言
	var (
		isCorrect int
	)

	// クエリを送信
	rows, _ := tx.Query(
		"SELECT COUNT(group_id) FROM users WHERE user_id = ? OR email = ? GROUP BY group_id;",
		userId,
		email,
	)
	defer rows.Close()
	
	for rows.Next() {
		rows.Scan(&isCorrect)
	}

	// 2人のグループIDが同じでない場合認証エラーを返す
	if isCorrect != 2 {
		return c.NoContent(http.StatusForbidden)
	}

	// 現在の管理者から管理者権限をなくす
	update, _ := tx.Prepare("UPDATE users SET is_manager = FALSE WHERE user_id = ?;")
	update.Exec(userId)

	// 新しい管理者に管理者権限を与える
	update, _ = tx.Prepare("UPDATE users SET is_manager = TRUE WHERE email = ?;")
	update.Exec(email)

	if err != nil {
		tx.Rollback()
		return c.NoContent(http.StatusInternalServerError)
	} else {
		tx.Commit()
		return c.NoContent(http.StatusOK)
	}
}