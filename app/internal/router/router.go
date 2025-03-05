package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"es-api/app/internal/handler"
	"es-api/app/middleware/cors"
)

func NewRouter(
	eh handler.ExperienceHandler,
	gh handler.ESGenerateHandler,
	authMiddleware echo.MiddlewareFunc,
) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	api := e.Group("/api")
	api.Use(cors.SetupCORS(e))
	api.Use(authMiddleware)
	api.GET("/experience", eh.GetExperienceByUserID)
	api.POST("/experience", eh.PostExperience)
	api.POST("/generate", gh.Generate)

	return e
}
