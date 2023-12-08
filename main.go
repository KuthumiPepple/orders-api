package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/kuthumipepple/orders-api/app"
)

func main() {
	application := app.New()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Fatal(application.Start(ctx))
}
