package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/nisarg1511/shortly/internal/handlers"
	"github.com/nisarg1511/shortly/internal/middleware"
	"github.com/nisarg1511/shortly/internal/models"
	"github.com/nisarg1511/shortly/internal/services"
	"github.com/nisarg1511/shortly/internal/store"
)

func main() {
	//Config Initialization
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg := models.Config{
		Address: ":" + port,
		DBSN:    os.Getenv("DB_CONN_STRING"),
	}
	db, err := sql.Open("pgx", cfg.DBSN)

	if err != nil {
		log.Fatalf("Critical: Couldn't parse DSN %v\n", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		log.Fatalf("Critical:Database is unreachable: %v\n", err)
	}

	log.Printf("Database connected successfully!\n")

	app := &models.Application{
		Config: cfg,
		DB:     db,
	}

	allStores := store.NewStore(app.DB)
	allServices := services.NewService(allStores)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         app.Config.Address,
		Handler:      middleware.RequestLogMiddleWare(middleware.RecoveryMiddleWare(mux)),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	//Services and handlers registration
	linkHandler := handlers.NewLink(allServices.Links)
	healthHandler := handlers.NewHealth(allServices.HealthService)
	//Health Routes
	mux.HandleFunc("GET /health", healthHandler.CheckHealth)
	//Link Routes
	mux.HandleFunc("POST /shorten", linkHandler.Shorten)
	mux.HandleFunc("GET /{code}", linkHandler.Redirect)

	//Server starting
	log.Printf("Server running on port:%s\n", port)
	log.Fatal(server.ListenAndServe())
}
