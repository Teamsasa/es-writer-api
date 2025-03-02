package main

import (
	"log"

	"github.com/joho/godotenv"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/handler"
	clerkRepo "es-api/app/internal/repository/clerk"
	dbRepo "es-api/app/internal/repository/db"
	"es-api/app/internal/router"
	"es-api/app/internal/usecase"
	"es-api/app/middleware/auth"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	dbConnManager := db.NewDBConnectionManager()
	experienceRepository := dbRepo.NewExperienceRepositoryWithDBManager(dbConnManager)
	clerkAuthRepository := clerkRepo.NewClerkAuthRepository()
	experienceUsecase := usecase.NewExperienceUsecase(experienceRepository)
	experienceHandler := handler.NewExperienceHandler(experienceUsecase)
	authMiddleware := auth.IDPAuthMiddleware(clerkAuthRepository, dbConnManager)
	e := router.NewRouter(experienceHandler, authMiddleware)
	e.Logger.Fatal(e.Start(":8080"))
}
