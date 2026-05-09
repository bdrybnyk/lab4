package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	dbURL := viper.GetString("db_url")
	port := viper.GetString("port")

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migrations applied successfully")

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	guitarRepo := NewGuitarRepo(dbpool)
	guitarHandler := NewGuitarHandler(guitarRepo)

	r := chi.NewRouter()

	r.Get("/guitars", guitarHandler.List)
	r.Get("/guitars/{id}", guitarHandler.GetById)
	r.Post("/guitars", guitarHandler.Create)
	r.Put("/guitars/{id}", guitarHandler.UpdateFull)
	r.Patch("/guitars/{id}", guitarHandler.UpdatePartial)
	r.Delete("/guitars/{id}", guitarHandler.Delete)

	log.Printf("Server is running on port %s\n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
