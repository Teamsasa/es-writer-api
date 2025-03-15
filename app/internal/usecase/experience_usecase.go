package usecase

import (
	"context"
	"fmt"

	"es-api/app/internal/entity/model"
	repository "es-api/app/internal/repository/db"
)

type ExperienceUsecase interface {
	GetExperienceByUserID(ctx context.Context) (*model.Experiences, error)
	PostExperience(ctx context.Context, experience model.InputExperience) (*model.Experiences, error)
}

type experienceUsecase struct {
	er repository.ExperienceRepository
}

func NewExperienceUsecase(r repository.ExperienceRepository) ExperienceUsecase {
	return &experienceUsecase{er: r}
}

func (u *experienceUsecase) GetExperienceByUserID(ctx context.Context) (*model.Experiences, error) {
	experience, err := u.er.GetExperienceByUserID(ctx)
	if err != nil {
		return nil, err
	}

	return &experience, nil
}

func (u *experienceUsecase) PostExperience(ctx context.Context, experience model.InputExperience) (*model.Experiences, error) {
	exists, err := u.er.FindExperienceByUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find experience by user ID: %w", err)
	}

	if exists {
		experiences, err := u.er.PatchExperience(ctx, experience)
		if err != nil {
			return nil, fmt.Errorf("failed to patch experience: %w", err)
		}

		return &experiences, nil
	}

	experiences, err := u.er.PostExperience(ctx, experience)
	if err != nil {
		return nil, fmt.Errorf("failed to post experience: %w", err)
	}

	return &experiences, nil
}
