package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"es-api/app/internal/contextKey"
	"es-api/app/internal/entity/model"
	repository "es-api/app/internal/repository/db"
	"es-api/app/test"
	"es-api/app/test/factory"
)

func TestExperienceRepository_GetExperienceByUserID(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	repo := repository.NewExperienceRepository(db)

	t.Run("異常系:ユーザーが存在しない場合", func(t *testing.T) {
		_ = factory.CreateUser2(t, db)

		ctx := test.SetupContextContext("test-user-id")
		ctx = context.WithValue(ctx, contextKey.UserIDKey, factory.DummyUserID1)
		experience, err := repo.GetExperienceByUserID(ctx)

		assert.Error(t, err)
		assert.Empty(t, experience)
	})

	t.Run("正常系:ユーザーが存在する場合", func(t *testing.T) {
		dummyUser := factory.CreateUser1(t, db)
		dummyExperience := factory.CreateExperience1(t, db)

		ctx := test.SetupContextContext("test-user-id")
		ctx = context.WithValue(ctx, contextKey.UserIDKey, dummyUser.ID)
		res, err := repo.GetExperienceByUserID(ctx)

		assert.NoError(t, err)
		assert.Equal(t, dummyExperience.ID, res.ID)
	})
}

func TestExperienceRepository_FindExperienceByUserID(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	repo := repository.NewExperienceRepository(db)
	dummyUser := factory.CreateUser1(t, db)

	t.Run("異常系:経験が存在しない場合", func(t *testing.T) {
		ctx := test.SetupContextContext("test-user-id")
		ctx = context.WithValue(ctx, contextKey.UserIDKey, dummyUser.ID)
		res, err := repo.FindExperienceByUserID(ctx)

		assert.NoError(t, err)
		assert.False(t, res)
	})

	t.Run("正常系:経験が存在する場合", func(t *testing.T) {
		_ = factory.CreateExperience1(t, db)

		ctx := test.SetupContextContext("test-user-id")
		ctx = context.WithValue(ctx, contextKey.UserIDKey, dummyUser.ID)
		res, err := repo.FindExperienceByUserID(ctx)

		assert.NoError(t, err)
		assert.True(t, res)
	})
}

func TestExperienceRepository_PostExperience(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	repo := repository.NewExperienceRepository(db)

	t.Run("正常系:経験を作成する", func(t *testing.T) {
		dummyUser := factory.CreateUser1(t, db)
		input := model.InputExperience{
			Work:        "test-work",
			Skills:      "test-skills",
			SelfPR:      "test-self-pr",
			FutureGoals: "test-future-goals",
		}

		ctx := test.SetupContextContext("test-user-id")
		ctx = context.WithValue(ctx, contextKey.UserIDKey, dummyUser.ID)
		experience, err := repo.PostExperience(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, experience.Work, input.Work)
		assert.Equal(t, experience.Skills, input.Skills)
		assert.Equal(t, experience.SelfPR, input.SelfPR)
		assert.Equal(t, experience.FutureGoals, input.FutureGoals)
	})
}

func TestExperienceRepository_PatchExperience(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	repo := repository.NewExperienceRepository(db)

	t.Run("正常系:経験を更新する", func(t *testing.T) {
		dummyUser := factory.CreateUser1(t, db)
		_ = factory.CreateExperience1(t, db)

		input := model.InputExperience{
			Work:        "updated-work",
			Skills:      "updated-skills",
			SelfPR:      "updated-self-pr",
			FutureGoals: "updated-future-goals",
		}

		ctx := test.SetupContextContext("test-user-id")
		ctx = context.WithValue(ctx, contextKey.UserIDKey, dummyUser.ID)
		experience, err := repo.PatchExperience(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, experience.Work, input.Work)
		assert.Equal(t, experience.Skills, input.Skills)
		assert.Equal(t, experience.SelfPR, input.SelfPR)
		assert.Equal(t, experience.FutureGoals, input.FutureGoals)
	})
}
