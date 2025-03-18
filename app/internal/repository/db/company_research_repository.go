package repository

import (
	"context"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/entity/model"
	"es-api/app/internal/contextKey"

	"gorm.io/gorm"
)

type CompanyResearchRepository interface {
	FindByCompanyID(ctx context.Context, companyID string) (*model.CompanyResearch, error)
	Create(ctx context.Context, research *model.CompanyResearch) error
}

type companyResearchRepository struct {
	dbManager db.DBConnectionManager
	defaultDB *gorm.DB
}

func NewCompanyResearchRepository(defaultDB *gorm.DB) CompanyResearchRepository {
	return &companyResearchRepository{
		defaultDB: defaultDB,
	}
}

func NewCompanyResearchRepositoryWithDBManager(dbManager db.DBConnectionManager) CompanyResearchRepository {
	return &companyResearchRepository{
		dbManager: dbManager,
		defaultDB: dbManager.GetConnection("clerk"),
	}
}

// FindByCompanyID - 法人番号で企業情報を検索
func (r *companyResearchRepository) FindByCompanyID(ctx context.Context, companyID string) (*model.CompanyResearch, error) {
	idp := ctx.Value(contextKey.IDPKey).(string)
	var dbConn *gorm.DB
	if r.dbManager != nil && idp != "" {
		dbConn = r.dbManager.GetConnection(idp)
	} else {
		dbConn = r.defaultDB
	}

	var research model.CompanyResearch
	result := dbConn.Where("company_id = ?", companyID).Find(&research)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &research, nil
}

// Create - 企業情報を新規作成
func (r *companyResearchRepository) Create(ctx context.Context, research *model.CompanyResearch) error {
	idp := ctx.Value(contextKey.IDPKey).(string)
	var dbConn *gorm.DB
	if r.dbManager != nil && idp != "" {
		dbConn = r.dbManager.GetConnection(idp)
	} else {
		dbConn = r.defaultDB
	}
	return dbConn.Create(research).Error
}
