package subject

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
)

type Subject struct {
	SubjectId	int
	SubjectName	string
}

type Subjects struct {
	Subjects	[]Subject
}

// func Register(c echo.Context) error {}

func Target(c echo.Context) error {
	classId   := c.FormValue("class_id")
	subjectId := c.FormValue("subject_id")

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO lesson_relations(class_id, subject_id) VALUES(?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	_, err = insert.Exec(classId, subjectId)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func LoadSubject(c echo.Context) error {
	classId := c.QueryParam("class_id")

	// データベースのハンドルを取得
	db := database.Connect()
	defer db.Close()

	rows, err := db.Query(
		"SELECT subject_id, subject_name FROM subjects NATURAL JOIN lesson_relations WHERE class_id = ?;",
		classId,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rows.Close()

	subjects := &Subjects{}
	for rows.Next() {
		subject := &Subject{}
		err := rows.Scan(&subject.SubjectId, &subject.SubjectName)

		if len(subjects.Subjects) == 0 {
			subjects.Subjects = []Subject{*subject}
		} else {
			subjects.Subjects = append(subjects.Subjects, *subject)
		}

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, subjects)
}