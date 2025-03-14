package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/usecase"
)

type ExperienceHandler interface {
	GetExperienceByUserID(c echo.Context) error
	PostExperience(c echo.Context) error
}

type experienceHandler struct {
	eu usecase.ExperienceUsecase
}

func NewExperienceHandler(eu usecase.ExperienceUsecase) ExperienceHandler {
	return &experienceHandler{eu: eu}
}

func (h *experienceHandler) GetExperienceByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	idp := c.Request().Header.Get("idp")
	userID := c.Get("userID")
	ctx = context.WithValue(ctx, "idp", idp)
	ctx = context.WithValue(ctx, "userID", userID)
	experience, err := h.eu.GetExperienceByUserID(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, experience)
}

func (h *experienceHandler) PostExperience(c echo.Context) error {
	var inputExperience model.InputExperience
	if err := c.Bind(&inputExperience); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	idp := c.Request().Header.Get("idp")
	userID := c.Get("userID")
	ctx = context.WithValue(ctx, "idp", idp)
	ctx = context.WithValue(ctx, "userID", userID)
	experience, err := h.eu.PostExperience(ctx, inputExperience)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, experience)
}
