package group

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {}

func Participate(c echo.Context) error {}

func Select(c echo.Context) error {}

func Quit(c echo.Context) error {}

func HostUpdate(c echo.Context) error {}