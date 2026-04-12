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

	err := http.ListenAndServe(addr, enableCors(rout))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func enableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
