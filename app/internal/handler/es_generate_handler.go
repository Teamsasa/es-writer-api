package handler

import (
	"log"
	"net/http"
	"strings"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/usecase"

	"github.com/labstack/echo/v4"
)

// ESGenerateHandler はES生成ハンドラーのインターフェース
type ESGenerateHandler interface {
	Generate(c echo.Context) error
}

type esGenerateHandler struct {
	esGenerateUsecase usecase.ESGenerateUsecase
}

// NewESGenerateHandler は新しいESGenerateHandlerを作成
func NewESGenerateHandler(esu usecase.ESGenerateUsecase) ESGenerateHandler {
	return &esGenerateHandler{
		esGenerateUsecase: esu,
	}
}

func (h *esGenerateHandler) Generate(c echo.Context) error {
	// リクエストをバインド
	req := new(model.ESGenerateRequest)
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

	// デフォルトモデル設定
	if req.Model == "" {
		req.Model = string(model.GeminiFlashLite)
	}

	// ユースケース内でユーザー情報の取得とバリデーションを行う
	result, err := h.esGenerateUsecase.GenerateES(c, *req)

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
