package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/database"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/handlers"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.DBConnect()

	repo := handlers.NewRepo(db)
	handlers.NewHandlers(repo)

	opts := []bot.Option{
		bot.WithDefaultHandler(repo.Handler),
	}

	b, err := bot.New(os.Getenv("TG_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	var mutex sync.Mutex
	go repo.TomorrowLeave(ctx, b, &mutex)

	b.Start(ctx)
}
