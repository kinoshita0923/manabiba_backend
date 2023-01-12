package class

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
)

type Class struct {
	ClassId		int
	ClassName	string
}

type GradeClasses struct {
	Grade		int
	Classes		[]Class
}

type GroupClasses struct {
	GroupName		string
	GradeClasses	[]GradeClasses
}

func Fetch(c echo.Context) error {
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusOK, "No token")
	}
	tokenText := tokenCookie.Value
	userId := jwt.ParseToken(tokenText).(float64)

	db := database.Connect()
	defer db.Close()

	var groupId int

	err = db.QueryRow(
		"SELECT group_id FROM users WHERE user_id = ?;",
		userId,
	).Scan(&groupId)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var (
		groupName string
		maxGrade  int
	)

	err = db.QueryRow(
		"SELECT group_name, max_grade FROM groups WHERE group_id = ?;",
		groupId,
	).Scan(&groupName, &maxGrade)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	groupClasses := &GroupClasses{GroupName: groupName}
	for i := 1; i <= maxGrade; i++ {
		gradeClasses := &GradeClasses{Grade: i}
		rows, err := db.Query(
			"SELECT class_id, class_name FROM classes WHERE grade = ?;",
			i,
		)
		if err != nil {
			log.Fatal(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		recordCount := 0
		for rows.Next() {
			recordCount++
			class := &Class{}
			err := rows.Scan(&class.ClassId, &class.ClassName)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			if len(gradeClasses.Classes) == 0 {
				gradeClasses.Classes = []Class{*class}
			} else {
				gradeClasses.Classes = append(gradeClasses.Classes, *class)
			}
		}

		if recordCount != 0 {
			if len(groupClasses.GradeClasses) == 0 {
				groupClasses.GradeClasses = []GradeClasses{*gradeClasses}
			} else {
				groupClasses.GradeClasses = append(groupClasses.GradeClasses, *gradeClasses)
			}
		}
	}

	return c.JSON(http.StatusOK, *groupClasses)
}

func Register(c echo.Context) error {
	className	:= c.FormValue("class_name")
	grade 		:= c.FormValue("grade")

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