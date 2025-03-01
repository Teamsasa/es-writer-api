package handler

import (
	"github.com/labstack/echo/v4"

	"tmp/app/internal/entity"
	"tmp/app/internal/usecase"
)

type WelcomeHandler interface {
	GetRandomUser(c echo.Context) error
}

type welcomeHandler struct {
	wu usecase.WelcomeUsecase
}

func NewWelcomeHandler(u usecase.WelcomeUsecase) WelcomeHandler {
	return &welcomeHandler{wu: u}
}

func (h *welcomeHandler) GetRandomUser(c echo.Context) error {
	var user *entity.User
	user, err := h.wu.GetRandomUser(c)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, user)
}
