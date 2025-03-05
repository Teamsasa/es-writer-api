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

// Generate godoc
// @Summary エントリーシートの回答を生成
// @Description HTMLから質問を抽出し、企業情報とユーザーの経験に基づいて回答を生成
// @Tags LLM
// @Accept json
// @Produce json
// @Param request body model.ESGenerateRequest true "ES生成リクエスト"
// @Success 200 {object} model.ESGenerateResponse "生成された回答"
// @Failure 400 {object} model.APIError "不正なリクエスト"
// @Failure 401 {object} model.APIError "認証エラー"
// @Failure 404 {object} model.APIError "リソースが見つからない"
// @Failure 500 {object} model.APIError "内部サーバーエラー"
// @Router /api/generate [post]
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
		switch {
		case isAuthError(err):
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
				"error": "認証に失敗しました",
			})
		case isNotFoundError(err):
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{
				"error": "リソースが見つかりません",
			})
		case isValidationError(err):
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
				"error": "回答生成中にエラーが発生しました",
			})
		}
	}

	return c.JSON(http.StatusOK, result)
}

// エラータイプ判定用のヘルパー関数
func isAuthError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "認証") ||
		strings.Contains(errorMsg, "ユーザーID") ||
		strings.Contains(errorMsg, "権限がありません")
}

func isNotFoundError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "見つかりません") ||
		strings.Contains(errorMsg, "存在しません")
}

func isValidationError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "無効な") ||
		strings.Contains(errorMsg, "必須です")
}
