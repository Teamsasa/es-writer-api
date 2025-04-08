package mock

import (
	"context"

	"es-api/app/internal/entity/model"

	"github.com/stretchr/testify/mock"
)

type GBizInfoRepositoryMock struct {
	mock.Mock
}

func (m *GBizInfoRepositoryMock) SearchCompanies(ctx context.Context, keyword string) ([]model.CompanyBasicInfo, error) {
	args := m.Called(ctx, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CompanyBasicInfo), args.Error(1)
}
