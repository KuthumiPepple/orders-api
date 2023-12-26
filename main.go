package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/kuthumipepple/orders-api/app"
	"github.com/kuthumipepple/orders-api/logger"
)

func main() {
	file, err := logger.SetupLog()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	application := app.New()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Fatal(application.Start(ctx))
}
