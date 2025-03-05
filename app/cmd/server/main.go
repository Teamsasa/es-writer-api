package main

import (
	"log"

	"github.com/joho/godotenv"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/handler"
	clerkRepo "es-api/app/internal/repository/clerk"
	dbRepo "es-api/app/internal/repository/db"
	geminiRepo "es-api/app/internal/repository/gemini"
	htmlRepo "es-api/app/internal/repository/parse_html"
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
	defaultDB := dbConnManager.GetConnection("default")

	// リポジトリの初期化
	experienceRepository := dbRepo.NewExperienceRepositoryWithDBManager(dbConnManager)
	authRepository := dbRepo.NewDBAuthRepository(defaultDB)
	clerkAuthRepository := clerkRepo.NewClerkAuthRepository()
	geminiRepository := geminiRepo.NewGeminiRepository()
	htmlAnalyzer := htmlRepo.NewHTMLAnalyzer(geminiRepository)
	tavilyRepository := tavilyRepo.NewTavilyRepository()

	// ユースケースの初期化
	experienceUsecase := usecase.NewExperienceUsecase(experienceRepository)
	esGenerateUsecase := usecase.NewESGenerateUsecase(
		htmlAnalyzer,
		geminiRepository,
		tavilyRepository,
		experienceRepository,
		authRepository,
	)

	// ハンドラーの初期化
	experienceHandler := handler.NewExperienceHandler(experienceUsecase)
	esGenerateHandler := handler.NewESGenerateHandler(esGenerateUsecase)

	// ミドルウェアの設定
	authMiddleware := auth.IDPAuthMiddleware(clerkAuthRepository, dbConnManager)

	// ルーターの設定とサーバーの起動
	e := router.NewRouter(experienceHandler, esGenerateHandler, authMiddleware)
	e.Logger.Fatal(e.Start(":8080"))
}
