package main

import (
	"fmt"
	"log"
	"net/http"

	"CPM_solver/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	rout := chi.NewRouter()

	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)

	rout.Post("/activity", handler.HandleActivity) // czynnosciowy wariant
	rout.Post("/event", handler.HandleEvent)       // zdarzeniowy

	addr := fmt.Sprintf(":%s", "8080")
	log.Printf("Listen on http://localhost%s", addr)

	err := http.ListenAndServe(addr, rout)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
