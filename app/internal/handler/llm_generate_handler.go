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
		log.Printf("LLM生成エラー: %v, 企業: %s", err, req.Company)

		if isAuthError(err) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "認証エラー",
			})
		}

		if isNotFoundError(err) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "リソースが見つかりません",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "回答生成中にエラーが発生しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"answers": result,
	})
}

func isAuthError(err error) bool {
	errMsg := err.Error()
	return contains(errMsg, []string{
		"認証",
		"unauthorized",
		"認証情報",
		"ログイン",
		"トークン",
		"セッション",
	})
}

func isNotFoundError(err error) bool {
	errMsg := err.Error()
	return contains(errMsg, []string{
		"見つかりません",
		"not found",
		"存在しません",
	})
}

func contains(s string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return false
}
