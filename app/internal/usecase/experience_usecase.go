package usecase

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"es-api/app/internal/entity/model"
	repository "es-api/app/internal/repository/db"
)

type ExperienceUsecase interface {
	GetExperienceByUserID(c echo.Context) (*model.Experiences, error)
	PostExperience(c echo.Context, experience model.InputExperience) (*model.Experiences, error)
}

type experienceUsecase struct {
	er repository.ExperienceRepository
}

func NewExperienceUsecase(r repository.ExperienceRepository) ExperienceUsecase {
	return &experienceUsecase{er: r}
}

func (u *experienceUsecase) GetExperienceByUserID(c echo.Context) (*model.Experiences, error) {
	experience, err := u.er.GetExperienceByUserID(c)
	if err != nil {
		return nil, err
	}

	return &experience, nil
}

func (u *experienceUsecase) PostExperience(c echo.Context, experience model.InputExperience) (*model.Experiences, error) {
	exists, err := u.er.FindExperienceByUserID(c)
	if err != nil {
		return nil, fmt.Errorf("failed to find experience by user ID: %w", err)
	}

	if exists {
		experiences, err := u.er.PatchExperience(c, experience)
		if err != nil {
			return nil, fmt.Errorf("failed to patch experience: %w", err)
		}

		return &experiences, nil
	}

	experiences, err := u.er.PostExperience(c, experience)
	if err != nil {
		return nil, fmt.Errorf("failed to post experience: %w", err)
	}

	return &experiences, nil
}
