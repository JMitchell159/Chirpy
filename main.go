package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JMitchell159/chirpy/internal/database"
	"github.com/JMitchell159/chirpy/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	apiCfg := &server.ApiConfig{}

	apiCfg.FileServerHits.Store(0)

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	apiCfg.DB = dbQueries

	apiCfg.Platform = os.Getenv("PLATFORM")

	server := apiCfg.CreateServer()
	server.ListenAndServe()
}
