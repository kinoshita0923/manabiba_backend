package main

import (
	"src/handlers/class"
	"src/handlers/evaluation"
	"src/handlers/examination"
	"src/handlers/group"
	"src/handlers/subject"
	"src/handlers/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
  // インスタンスを作成
  e := echo.New()

  // ミドルウェアを設定
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
  }))

  // 静的サイトを返すルーティングの設定
  e.Static("/docs/", "../docs/")

  e.Static("/", "./frontend/dist")
  e.Static("/login", "./frontend/dist")
  e.Static("/signup", "./frontend/dist")
  e.Static("/search-group", "./frontend/dist")

  // APIのルーティング
  e.POST("/user/register", user.Register)
  e.POST("/user/authentication", user.Authentication)
  e.GET("/user/check-login", user.CheckLogin)
  e.POST("/group/register", group.Register)
  e.POST("/group/participate", group.Participate)
  e.GET("/group/select", group.Select)
  e.PATCH("/group/quit", group.Quit)
  e.PATCH("/group/host-update", group.HostUpdate)
  // e.POST("/subject/register", subject.Register)
  e.GET("/subject/load-subject", subject.LoadSubject)
  e.GET("/evaluation/fetch", evaluation.Fetch)
  e.POST("/evaluation/register", evaluation.Register)
  e.POST("/evaluation/purchase", evaluation.Purchase)
  e.POST("/examination/upload", examination.Upload)
  e.POST("/examination/purchase", examination.Purchase)
  e.GET("/class/fetch", class.Fetch)
  e.POST("/class/register", class.Register)
  // e.PUT("/good", good.Reverse)

  // サーバーをポート番号8080で起動
  e.Logger.Fatal(e.Start(":8080"))
}
