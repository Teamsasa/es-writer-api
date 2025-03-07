package handler

import (
	"log"
	"net/http"
	"strings"

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
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "リクエストの解析に失敗しました",
		})
	}

	// 基本的なバリデーション
	if req.Company == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "企業名は必須です",
		})
	}
	if req.HTML == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "HTMLは必須です",
		})
	}


	// ユースケース内でユーザー情報の取得とバリデーションを行う
	result, err := h.llmenerateUsecase.LLMGenerate(c, *req)

	// エラーハンドリング
	if err != nil {
		log.Printf("ES生成エラー: %v, 企業: %s", err, req.Company)

		// エラータイプに応じたHTTPステータス
		if isValidationError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
				"error": "回答生成中にエラーが発生しました",
			})
		}
	}

	return c.JSON(http.StatusOK, result)
}

func isValidationError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "無効な") ||
		strings.Contains(errorMsg, "必須です")
}
