package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"

	"es-api/app/infrastructure/db"
	clerkRepo "es-api/app/internal/repository/clerk"
	dbRepo "es-api/app/internal/repository/db"
)

func IDPAuthMiddleware(
	clerkAuthRepo clerkRepo.ClerkAuthRepository,
	dbConnManager db.DBConnectionManager,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idp := c.Request().Header.Get("idp")

			if idp == "swagger" || idp == "test" {
				dbConn := dbConnManager.GetConnection(idp)

				dummyUserID := "user_abcdefghijklmnopqrstuvwxyza"
				if idp == "test" {
					dummyUserID = "test-user"
				}

				c.Set("userID", dummyUserID)

				dbAuthRepo := dbRepo.NewDBAuthRepository(dbConn)

				exists, err := dbAuthRepo.FindUser(dummyUserID)
				if err != nil {
					fmt.Printf("Failed to check user existence: %v\n", err)
				} else if !exists {
					err = dbAuthRepo.CreateUser(dummyUserID)
					if err != nil {
						fmt.Printf("Failed to create user: %v\n", err)
					} else {
						fmt.Printf("Created new dummy user with ID: %s\n", dummyUserID)
					}
				}
				return next(c)
			}
			return clerkAuthentication(c, next, clerkAuthRepo, dbConnManager)
		}
	}
}

func clerkAuthentication(
	c echo.Context,
	next echo.HandlerFunc,
	clerkAuthRepo clerkRepo.ClerkAuthRepository,
	dbConnManager db.DBConnectionManager,
) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authentication format"})
	}

	tokenString := parts[1]

	keySet, err := clerkAuthRepo.FetchJWKS()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch JWKS: %v", err),
		})
	}

	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": fmt.Sprintf("Token validation failed: %v", err),
		})
	}

	userID, ok := token.Get("sub")
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Sub claim not found in token",
		})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Sub claim is not a string",
		})
	}

	c.Set("userID", userIDStr)

	dbConn := dbConnManager.GetConnection("clerk")
	dbAuthRepo := dbRepo.NewDBAuthRepository(dbConn)

	exists, err := dbAuthRepo.FindUser(userIDStr)
	if err != nil {
		fmt.Printf("Failed to check user existence: %v\n", err)
	} else if !exists {
		err = dbAuthRepo.CreateUser(userIDStr)
		if err != nil {
			fmt.Printf("Failed to create user: %v\n", err)
		} else {
			fmt.Printf("Created new user with ID: %s\n", userIDStr)
		}
	}

	return next(c)
}
