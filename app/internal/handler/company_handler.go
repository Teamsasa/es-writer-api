package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"es-api/app/internal/usecase"
)

type CompanyHandler interface {
	SearchCompanies(c echo.Context) error
}

type companyHandler struct {
	companyUsecase usecase.CompanyUsecase
}

func NewCompanyHandler(companyUsecase usecase.CompanyUsecase) CompanyHandler {
	return &companyHandler{
		companyUsecase: companyUsecase,
	}
}

func (h *companyHandler) SearchCompanies(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "keyword is required",
		})
	}

	companies, err := h.companyUsecase.SearchCompanies(c, keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to search companies: %v", err),
		})
	}

	return c.JSON(http.StatusOK, companies)
}
