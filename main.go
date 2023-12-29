package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/kuthumipepple/orders-api/app"
	"github.com/kuthumipepple/orders-api/internal/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("unable to load `.env` file:", err)
	}

	file, err := logger.SetupLog()
	if err != nil {
		log.Println("failed to set up log file:", err)
	}
	defer file.Close()

	application := app.New(app.LoadConfig())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Fatal(application.Start(ctx))
}
