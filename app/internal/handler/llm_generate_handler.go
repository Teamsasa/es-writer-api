package handler

import (
	"log"
	"net/http"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/usecase"

	"github.com/labstack/echo/v4"
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

	if req.Company == "" || req.HTML == "" {
		log.Printf("必須パラメータ不足: company=%v, html=%v", req.Company != "", req.HTML != "")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "必要なパラメータが不足しています",
		})
	}

	result, err := h.llmenerateUsecase.LLMGenerate(c, *req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),	
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"answers": result,
	})
}
