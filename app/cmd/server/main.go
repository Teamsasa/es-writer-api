package main

import (
	"log"

	"github.com/joho/godotenv"

	"tmp/app/infrastructure/db"
	"tmp/app/internal/handler"
	"tmp/app/internal/repository"
	"tmp/app/internal/router"
	"tmp/app/internal/usecase"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	dbConnection := db.NewDB()
	welcomeRepository := repository.NewWelcomeRepository(dbConnection)
	welcomeUsecase := usecase.NewWelcomeUsecase(welcomeRepository)
	welcomeHandler := handler.NewWelcomeHandler(welcomeUsecase)
	e := router.NewRouter(welcomeHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
