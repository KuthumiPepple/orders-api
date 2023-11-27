package app

import (
	"context"
	"fmt"
	"net/http"
)

type Application struct {
	router http.Handler
}

func New() *Application {
	return &Application{
		router: loadRoutes(),
	}
}

func (a *Application) Start(c context.Context) error {
	err := http.ListenAndServe(":8000", a.router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
