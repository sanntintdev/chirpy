package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/handlers"
	"github.com/sanntintdev/chirpy/pkg/middleware"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	dbQueries := database.New(db)
	cfg := config.NewAPIConfig(dbQueries)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// Static file serving
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", middleware.MetricsInc(cfg, http.FileServer(http.Dir("./assets")))))
	mux.Handle("/app/", http.StripPrefix("/app/", middleware.MetricsInc(cfg, http.FileServer(http.Dir(".")))))

	// Admin routes
	mux.HandleFunc("GET /admin/metrics", handlers.MetricsHandler(cfg))
	mux.HandleFunc("POST /admin/reset", handlers.ResetHandler(cfg))

	// API routes
	mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler(cfg))
	mux.HandleFunc("POST /api/users", handlers.CreateUserHandler(cfg))
	mux.HandleFunc("GET /api/chirps", handlers.GetChirpsHandler(cfg))
	mux.Handle("POST /api/chirps", handlers.CreateChirpHandler(cfg))
	mux.Handle("GET /api/chirps/{chirpId}", handlers.GetChirpHandler(cfg))

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
