package handler

import (
	"context"
	"log"
	"net/http"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/usecase"

	"github.com/labstack/echo/v4"

	"es-api/app/internal/contextKey"
)

type LLMGenerateHandler interface {
	Generate(c echo.Context) error
}

type llmGenerateHandler struct {
	llmenerateUsecase usecase.LLMGenerateUsecase
}

func NewLLMGenerateHandler(llmu usecase.LLMGenerateUsecase) LLMGenerateHandler {
	return &llmGenerateHandler{
		llmenerateUsecase: llmu,
	}
}

func (h *llmGenerateHandler) Generate(c echo.Context) error {
	// リクエストをバインド
	req := new(model.LLMGenerateRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("リクエストバインドエラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "リクエストの解析に失敗しました",
		})
	}

	if req.CompanyName == "" || req.CompanyID == "" || req.HTML == "" {
		log.Printf("必須パラメータ不足: companyName=%v, companyID=%v, html=%v", req.CompanyName != "", req.CompanyID != "", req.HTML != "")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "必要なパラメータが不足しています",
		})
	}

	ctx := c.Request().Context()
	idp := c.Request().Header.Get("idp")
	userID := c.Get("userID")
	ctx = context.WithValue(ctx, contextKey.IDPKey, idp)
	ctx = context.WithValue(ctx, contextKey.UserIDKey, userID)

	result, err := h.llmenerateUsecase.LLMGenerate(ctx, *req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"answers": result,
	})
}
