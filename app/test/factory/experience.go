package factory

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"es-api/app/internal/entity/model"
)

const (
	DummyExperienceID1 = "123e4567-e89b-12d3-a456-426614174000"
	DummyExperienceID2 = "123e4567-e89b-12d3-a456-426614174001"
)

func CreateExperience1(t *testing.T, dbConn *gorm.DB) model.Experiences {
	experience := model.Experiences{
		ID:          DummyExperienceID1,
		UserID:      DummyUserID1,
		Work:        "work",
		Skills:      "skills",
		SelfPR:      "selfPR",
		FutureGoals: "futureGoals",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := dbConn.Create(&experience).Error; err != nil {
		t.Fatalf("Error creating test experience: %v", err)
	}

	return experience
}

func CreateExperience2(t *testing.T, dbConn *gorm.DB) model.Experiences {
	experience := model.Experiences{
		ID:          DummyExperienceID2,
		UserID:      DummyUserID2,
		Work:        "work2",
		Skills:      "skills2",
		SelfPR:      "selfPR2",
		FutureGoals: "futureGoals2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := dbConn.Create(&experience).Error; err != nil {
		t.Fatalf("Error creating test experience: %v", err)
	}

	return experience
}
