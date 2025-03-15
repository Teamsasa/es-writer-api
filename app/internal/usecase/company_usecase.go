package usecase

import (
	"context"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/repository/gbiz"
)

type CompanyUsecase interface {
	SearchCompanies(ctx context.Context, keyword string) ([]model.CompanyBasicInfo, error)
}

type companyUsecase struct {
	gbizRepo gbiz.GBizInfoRepository
}

func NewCompanyUsecase(gbizRepo gbiz.GBizInfoRepository) CompanyUsecase {
	return &companyUsecase{
		gbizRepo: gbizRepo,
	}
}

func (u *companyUsecase) SearchCompanies(ctx context.Context, keyword string) ([]model.CompanyBasicInfo, error) {
	return u.gbizRepo.SearchCompanies(ctx, keyword)
}
