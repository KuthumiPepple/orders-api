package main

import (
	"context"
	"log"

	"github.com/kuthumipepple/orders-api/app"
)

func main() {
	application := app.New()
	log.Fatal(application.Start(context.TODO()))
}
