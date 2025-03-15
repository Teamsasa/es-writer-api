package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/entity/model"
)

type ExperienceRepository interface {
	GetExperienceByUserID(ctx context.Context) (model.Experiences, error)
	FindExperienceByUserID(ctx context.Context) (bool, error)
	PostExperience(ctx context.Context, input model.InputExperience) (model.Experiences, error)
	PatchExperience(ctx context.Context, input model.InputExperience) (model.Experiences, error)
}

type experienceRepository struct {
	dbManager db.DBConnectionManager
	defaultDB *gorm.DB
}

func NewExperienceRepository(defaultDB *gorm.DB) ExperienceRepository {
	return &experienceRepository{
		defaultDB: defaultDB,
	}
}

func NewExperienceRepositoryWithDBManager(dbManager db.DBConnectionManager) ExperienceRepository {
	return &experienceRepository{
		dbManager: dbManager,
		defaultDB: dbManager.GetConnection("clerk"),
	}
}

func (r *experienceRepository) GetExperienceByUserID(ctx context.Context) (model.Experiences, error) {
	var experience model.Experiences
	idp := ctx.Value("idp").(string)
	userID := ctx.Value("userID").(string)
	var dbConn *gorm.DB
	if r.dbManager != nil && idp != "" {
		dbConn = r.dbManager.GetConnection(idp)
	} else {
		dbConn = r.defaultDB
	}
	result := dbConn.First(&experience, "user_id = ?", userID)
	if result.Error != nil {
		return model.Experiences{}, result.Error
	}
	return experience, nil
}

func (r *experienceRepository) FindExperienceByUserID(ctx context.Context) (bool, error) {
	var experience model.Experiences
	idp := ctx.Value("idp").(string)
	userID := ctx.Value("userID").(string)
	var dbConn *gorm.DB
	if r.dbManager != nil && idp != "" {
		dbConn = r.dbManager.GetConnection(idp)
	} else {
		dbConn = r.defaultDB
	}
	result := dbConn.Where("user_id = ?", userID).First(&experience)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (r *experienceRepository) PostExperience(ctx context.Context, input model.InputExperience) (model.Experiences, error) {
	idp := ctx.Value("idp").(string)
	userID := ctx.Value("userID").(string)
	var dbConn *gorm.DB
	if r.dbManager != nil && idp != "" {
		dbConn = r.dbManager.GetConnection(idp)
	} else {
		dbConn = r.defaultDB
	}
	experience := model.Experiences{
		UserID:      userID,
		Work:        input.Work,
		Skills:      input.Skills,
		SelfPR:      input.SelfPR,
		FutureGoals: input.FutureGoals,
	}

	result := dbConn.Create(&experience)
	if result.Error != nil {
		return model.Experiences{}, result.Error
	}
	return experience, nil
}

func (r *experienceRepository) PatchExperience(ctx context.Context, input model.InputExperience) (model.Experiences, error) {
	idp := ctx.Value("idp").(string)
	userID := ctx.Value("userID").(string)
	var dbConn *gorm.DB
	if r.dbManager != nil && idp != "" {
		dbConn = r.dbManager.GetConnection(idp)
	} else {
		dbConn = r.defaultDB
	}

	var experience model.Experiences
	_ = dbConn.Where("user_id = ?", userID).First(&experience)

	experience.Work = input.Work
	experience.Skills = input.Skills
	experience.SelfPR = input.SelfPR
	experience.FutureGoals = input.FutureGoals

	result := dbConn.Save(&experience)
	if result.Error != nil {
		return model.Experiences{}, result.Error
	}
	return experience, nil
}
