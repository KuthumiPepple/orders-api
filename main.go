package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", basicHandler)

	http.ListenAndServe(":8000", router)
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("basic server is running!"))
}
