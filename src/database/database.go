package database

import (
	"fmt"
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	// データベースの情報を定義
	_ = godotenv.Load(".env")
	USER 	 := os.Getenv("USER_NAME")
	PASSWORD := os.Getenv("PASSWORD")
	HOST	 := os.Getenv("HOST")
	PORT 	 := os.Getenv("PORT")
	DATABASE := os.Getenv("DATABASE")

	// データベースのハンドルを取得する
	dbconf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USER, PASSWORD, HOST, PORT, DATABASE)
	db, _ := sql.Open("mysql", dbconf)

	return db
}
