package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
  // インスタンスを作成
  e := echo.New()

  // ミドルウェアを設定
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // 静的サイトを返すルーティングの設定
  e.Static("/", "../../frontend/dist/")

  // サーバーをポート番号80で起動
  e.Logger.Fatal(e.Start(":80"))
}
