package usecase

import (
	"es-api/app/internal/entity/model"
	"es-api/app/internal/repository/gbiz"

	"github.com/labstack/echo/v4"
)

type CompanyUsecase interface {
	SearchCompanies(c echo.Context, keyword string) ([]model.CompanyBasicInfo, error)
}

type companyUsecase struct {
	gbizRepo gbiz.GBizInfoRepository
}

func NewCompanyUsecase(gbizRepo gbiz.GBizInfoRepository) CompanyUsecase {
	return &companyUsecase{
		gbizRepo: gbizRepo,
	}
}

func (u *companyUsecase) SearchCompanies(c echo.Context, keyword string) ([]model.CompanyBasicInfo, error) {
	return u.gbizRepo.SearchCompanies(c, keyword)
}
