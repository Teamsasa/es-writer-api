package main

import (
	"log"

	"github.com/joho/godotenv"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/handler"
	clerkRepo "es-api/app/internal/repository/clerk"
	dbRepo "es-api/app/internal/repository/db"
	geminiRepo "es-api/app/internal/repository/gemini"
	tavilyRepo "es-api/app/internal/repository/tavily"
	"es-api/app/internal/router"
	"es-api/app/internal/usecase"
	"es-api/app/middleware/auth"
)

func main() {
	// 環境変数の読み込み
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}

	// データベース接続の初期化
	dbConnManager := db.NewDBConnectionManager()

	// リポジトリの初期化
	experienceRepository := dbRepo.NewExperienceRepositoryWithDBManager(dbConnManager)
	clerkAuthRepository := clerkRepo.NewClerkAuthRepository()
	geminiRepository := geminiRepo.NewGeminiRepository()
	tavilyRepository := tavilyRepo.NewTavilyRepository()
	htmlExtractUsecase := usecase.NewHTMLExtractUsecase(geminiRepository)

	// ユースケースの初期化
	experienceUsecase := usecase.NewExperienceUsecase(experienceRepository)
	llmGenerateUsecase := usecase.NewLLMGenerateUsecase(
		htmlExtractUsecase,
		geminiRepository,
		tavilyRepository,
		experienceRepository,
		nil,
	)

	// ハンドラーの初期化
	experienceHandler := handler.NewExperienceHandler(experienceUsecase)
	llmGenerateHandler := handler.NewLLMGenerateHandler(llmGenerateUsecase)

	// ミドルウェアの設定
	authMiddleware := auth.IDPAuthMiddleware(clerkAuthRepository, dbConnManager)

	// ルーターの設定とサーバーの起動
	e := router.NewRouter(experienceHandler, llmGenerateHandler, authMiddleware)
	e.Logger.Fatal(e.Start(":8080"))
}
