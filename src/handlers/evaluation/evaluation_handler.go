package evaluation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
)

type Evaluation struct {
	EvaluationId	int
	Evaluation		int
	Comment			string
	TeacherName		string
	Term			int
	StudyTime		float64
}

type Evaluations struct {
	Evaluations []Evaluation
}

func Fetch(c echo.Context) error {
	subjectId := c.QueryParam("subject_id")

	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusOK, "No token")
	}
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	// DBにクエリを送信
	rows, err := db.Query(
		"SELECT	evaluation_id, evaluation, comment, teacher_name, term, study_time FROM evaluations NATURAL JOIN viewable_contents WHERE subject_id = ? AND user_id = ?	AND confirm_genre = FALSE;",
		subjectId,
		userId,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rows.Close()

	// 結果を代入
	evaluations := &Evaluations{}
	for rows.Next() {
		evaluation := &Evaluation{}
		err := rows.Scan(&evaluation.EvaluationId, &evaluation.Evaluation, &evaluation.Comment, &evaluation.TeacherName, &evaluation.Term, &evaluation.StudyTime)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if len(evaluations.Evaluations) == 0 {
			evaluations.Evaluations = []Evaluation{*evaluation}
		} else {
			evaluations.Evaluations = append(evaluations.Evaluations, *evaluation)
		}
	}

	fmt.Println(evaluations)

	return c.JSON(http.StatusOK, evaluations)
}

func Register(c echo.Context) error {
	subjectId	:= c.FormValue("subject_id")
	evaluation	:= c.FormValue("evaluation")
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
	insert, err := db.Prepare("INSERT INTO evaluations(subject_id, user_id, evaluation, comment, teacher_name, term, study_time) VALUES(?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	_, err = insert.Exec(subjectId, userId, evaluation, comment, teacherName, term, studyTime)
	if err != nil{
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func Purchase(c echo.Context) error {
	subjectId := c.FormValue("subject_id")

	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusOK, "No token")
	}
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO viewable_contents(user_id, subject_id, confirm_genre) VALUES(?, ?, FALSE);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	_, err = insert.Exec(userId, subjectId)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}
