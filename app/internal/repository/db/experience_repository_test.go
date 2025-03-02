package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"es-api/app/internal/entity/model"
	repository "es-api/app/internal/repository/db"
	"es-api/app/test"
	"es-api/app/test/factory"
)

// (GetExperienceByUserID - 正常系　- ユーザーが存在する)
func TestExperienceRepository_GetExperienceByUserID(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	dummyUser := factory.CreateUser1(t, db)
	dummyExperience := factory.CreateExperience1(t, db)
	repo := repository.NewExperienceRepository(db)

	ctx := test.SetupEchoContext(dummyUser.ID)
	res, err := repo.GetExperienceByUserID(ctx)

	assert.NoError(t, err)
	assert.Equal(t, dummyExperience.ID, res.ID)
}

// (GetExperienceByUserID - 異常系　- ユーザーが存在しない)
func TestExperienceRepository_GetExperienceByUserID_NotFound(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	_ = factory.CreateUser2(t, db)
	repo := repository.NewExperienceRepository(db)

	ctx := test.SetupEchoContext(factory.DummyUserID1)
	experience, err := repo.GetExperienceByUserID(ctx)

	assert.Error(t, err)
	assert.Empty(t, experience)
}

// (FindExperienceByUserID - 正常系　- 経験が存在する)
func TestExperienceRepository_FindExperienceByUserID_Found(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	dummyUser := factory.CreateUser1(t, db)
	_ = factory.CreateExperience1(t, db)
	repo := repository.NewExperienceRepository(db)

	ctx := test.SetupEchoContext(dummyUser.ID)
	res, _ := repo.FindExperienceByUserID(ctx)

	assert.True(t, res)
}

// (FindExperienceByUserID - 異常系　- 経験が存在しない)
func TestExperienceRepository_FindExperienceByUserID_NotFound(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	dummyUser := factory.CreateUser1(t, db)
	repo := repository.NewExperienceRepository(db)

	ctx := test.SetupEchoContext(dummyUser.ID)
	res, _ := repo.FindExperienceByUserID(ctx)

	assert.False(t, res)
}

// (PostExperience - 正常系　- 経験を作成する)
func TestExperienceRepository_PostExperience(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	dummyUser := factory.CreateUser1(t, db)
	repo := repository.NewExperienceRepository(db)
	input := model.InputExperience{
		Work:        "test-work",
		Skills:      "test-skills",
		SelfPR:      "test-self-pr",
		FutureGoals: "test-future-goals",
	}

	ctx := test.SetupEchoContext(dummyUser.ID)
	experience, err := repo.PostExperience(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, experience.Work, input.Work)
	assert.Equal(t, experience.Skills, input.Skills)
	assert.Equal(t, experience.SelfPR, input.SelfPR)
	assert.Equal(t, experience.FutureGoals, input.FutureGoals)
}
