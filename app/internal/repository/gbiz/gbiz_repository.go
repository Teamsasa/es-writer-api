package gbiz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"es-api/app/internal/entity/model"

	"github.com/labstack/echo/v4"
)

type GBizInfoRepository interface {
	SearchCompanies(c echo.Context, keyword string) ([]model.CompanyBasicInfo, error)
}

type gbizInfoRepository struct {
	baseURL string
}

func NewGBizInfoRepository() GBizInfoRepository {
	return &gbizInfoRepository{
		baseURL: "https://info.gbiz.go.jp/hojin/v1/hojin",
	}
}

// SearchCompanies - 法人名の検索を行う
func (r *gbizInfoRepository) SearchCompanies(c echo.Context, keyword string) ([]model.CompanyBasicInfo, error) {
	apiKey := os.Getenv("GBIZ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GBIZ_API_KEY is not set")
	}

	// リクエストの構築
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?name=%s&exist_flg=true", r.baseURL, url.QueryEscape(keyword)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// APIキーをヘッダーに設定
	req.Header.Set("X-hojinInfo-api-token", apiKey)

	// リクエストの実行
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// エラーレスポンスの確認
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// レスポンスのパース
	var gbizResponse model.GBizInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&gbizResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 結果の変換
	companies := make([]model.CompanyBasicInfo, 0, len(gbizResponse.Response))
	for _, company := range gbizResponse.Response {
		companies = append(companies, model.CompanyBasicInfo{
			CompanyID:   company.CorporateNumber,
			CompanyName: company.Name,
		})
	}

	return companies, nil
}
