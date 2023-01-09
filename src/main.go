package main

import (
  "src/handlers/user"
  "src/handlers/group"

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

  e.Static("/", "../frontend/dist")
  e.Static("/login", "../frontend/dist")
  e.Static("/signup", "../frontend/dist")
  e.Static("/search-group", "../frontend/dist")

  // APIのルーティング
  e.POST("/user/register", user.Register)
  e.POST("/user/authentication", user.Authentication)
  e.GET("/user/check-login", user.CheckLogin)
  e.POST("/group/register", group.Register)
  e.POST("/group/participate", group.Participate)
  // e.GET("/group/select", group.Select)
  // e.PATCH("/group/quit", group.Quit)
  // e.PATCH("/group/host-update", group.HostUpdate)
  // e.POST("/subject/register", subject.Register)
  // e.GET("/subject/select", subject.Select)
  // e.POST("/point/add-difference", point.AddDifference)
  // e.POST("/content/comment/register", comment.Register)
  // e.PATCH("/content/comment/update", comment.Update)
  // e.POST("/content/examination/register", examination.Register)
  // e.PATCH("/content/examination/update", examination.Update)
  // e.GET("/content/select", content.Select)
  // e.DELETE("/content/delete", content.Delete)
  // e.PUT("/good", good.Reverse)

  // サーバーをポート番号8080で起動
  e.Logger.Fatal(e.Start(":8080"))
}
