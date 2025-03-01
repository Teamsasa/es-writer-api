package repository_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"tmp/app/internal/repository"
	"tmp/app/test"
)

func TestWelcomeRepository_GetAllUsers(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupDB(t, db)

	expectedUser := test.SeedTestUser(t, db)
	repo := repository.NewWelcomeRepository(db)

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	users, err := repo.GetAllUsers(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, expectedUser.ID, users[0].ID)
	assert.Equal(t, expectedUser.Name, users[0].Name)
	assert.Equal(t, expectedUser.Email, users[0].Email)
}

func TestWelcomeRepository_GetAllUsers_Empty(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupDB(t, db)

	repo := repository.NewWelcomeRepository(db)

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	users, err := repo.GetAllUsers(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
