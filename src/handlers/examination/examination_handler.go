package examination

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"src/database"
	"src/lib/jwt"
)

// 過去問をアップロードする関数
func Upload(c echo.Context) error {
	subjectId   := c.FormValue("subject_id")
	nth			:= c.FormValue("nth_quarter")
	studyTime	:= c.FormValue("study_time")

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
	insert, err := db.Prepare("INSERT INTO examinations(user_id, subject_id, nth_quarter, study_time) VALUES(?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer insert.Close()

	// SQLの実行
	res, err := insert.Exec(userId, subjectId, nth, studyTime)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// 過去問IDを取得
	examinationId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// 問題と模範解答を取得
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	quizFiles   := form.File["quizes"]
	answerFiles := form.File["answers"]

	fileId := 1

	// 過去問の数だけ画像を保存
	for _, file := range quizFiles {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		fileName := fmt.Sprintf("storage/%d_%d.jpg", examinationId, fileId)

		dst, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// 画像のパスをデータベースに保存
		insert, err := db.Prepare("INSERT INTO image_paths(examination_id, file_id, path, is_answer) VALUES(?, ?, ?, FALSE);")
		if err != nil {
			log.Fatal(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer insert.Close()

		_, err = insert.Exec(examinationId, fileId, fileName)
		if err != nil {
			log.Fatal(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		fileId++
	}

	// 模範解答の数だけ画像を保存
	for _, file := range answerFiles {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		fileName := fmt.Sprintf("storage/%d_%d.jpg", examinationId, fileId)

		dst, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// 画像のパスをデータベースに保存
		insert, err := db.Prepare("INSERT INTO image_paths(examination_id, file_id, path, is_answer) VALUES(?, ?, ?, TRUE);")
		if err != nil {
			log.Fatal(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer insert.Close()

		_, err = insert.Exec(examinationId, fileId, fileName)
		if err != nil {
			log.Fatal(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		fileId++
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

	insert, err := db.Prepare("INSERT INTO viewable_contents(user_id, subject_id, confirm_genre) VALUES(?, ?, TRUE);")
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