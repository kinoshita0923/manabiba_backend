package evaluation

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
)

func Register(c echo.Context) error {
	subjectId	:= c.FormValue("subject_id")
	valuation	:= c.FormValue("valuation")
	comment		:= c.FormValue("comment")
	teacherName	:= c.FormValue("teacher_name")
	term		:= c.FormValue("term")
	studyTime	:= c.FormValue("study_time")

	// ユーザIDをトークンから取得
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusOK, "No token")
	}
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// SQL文を定義
	insert, err := db.Prepare("INSERT INTO evaluations(subject_id, user_id, valuation, comment, teacher_name, term, study_time) VALUES(?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	_, err = insert.Exec(subjectId, userId, valuation, comment, teacherName, term, studyTime)
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

// func Update(c echo.Context) error {}