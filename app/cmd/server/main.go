package main

import (
	"log"

	"github.com/joho/godotenv"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/handler"
	clerkRepo "es-api/app/internal/repository/clerk"
	dbRepo "es-api/app/internal/repository/db"
	gbizRepo "es-api/app/internal/repository/gbiz"
	geminiRepo "es-api/app/internal/repository/gemini"
	tavilyRepo "es-api/app/internal/repository/tavily"
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
	companyResearchRepository := dbRepo.NewCompanyResearchRepositoryWithDBManager(dbConnManager)
	clerkAuthRepository := clerkRepo.NewClerkAuthRepository()
	geminiRepository := geminiRepo.NewGeminiRepository()
	tavilyRepository := tavilyRepo.NewTavilyRepository()
	gbizRepository := gbizRepo.NewGBizInfoRepository()
	experienceUsecase := usecase.NewExperienceUsecase(experienceRepository)
	companyUsecase := usecase.NewCompanyUsecase(gbizRepository)
	llmGenerateUsecase := usecase.NewLLMGenerateUsecase(
		geminiRepository,
		tavilyRepository,
		experienceRepository,
		companyResearchRepository,
	)
	experienceHandler := handler.NewExperienceHandler(experienceUsecase)
	llmGenerateHandler := handler.NewLLMGenerateHandler(llmGenerateUsecase)
	companyHandler := handler.NewCompanyHandler(companyUsecase)
	authMiddleware := auth.IDPAuthMiddleware(clerkAuthRepository, dbConnManager)
	e := router.NewRouter(experienceHandler, llmGenerateHandler, companyHandler, authMiddleware)
	e.Logger.Fatal(e.Start(":8080"))
}
