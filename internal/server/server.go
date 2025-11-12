package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/thehappyidiot/save-my-music/internal/database"
)

type Config struct {
	port         int
	isProduction bool
}

type Server struct {
	config    Config
	dbQueries *database.Queries
}

func NewServer() *http.Server {
	// Get config:
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Cannot parse environment variable `port` as int")
	}

	isProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	if err != nil {
		log.Fatal("Cannot parse environment variable `is_production` as boolean")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to database: %s", err))
	}

	server := Server{
		config: Config{
			port:         port,
			isProduction: isProduction,
		},
		dbQueries: database.New(db),
	}

	httpServer := &http.Server{
		Handler: server.RegisterRoutes(),
		Addr:    fmt.Sprintf(":%d", port),
	}

	return httpServer
}
