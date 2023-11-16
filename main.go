package main

import (
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(basicHandler),
	}

	log.Fatal(server.ListenAndServe())
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("basic server is running!"))
}	
