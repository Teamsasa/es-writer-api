package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"tmp/app/internal/handler"
)

func NewRouter(
	wh handler.WelcomeHandler,
) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.GET("/", wh.GetRandomUser)

	return e
}
