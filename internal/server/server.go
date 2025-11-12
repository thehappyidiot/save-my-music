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
	port          int
	isDevelopment bool
}

type Server struct {
	config         Config
	dbQueries      *database.Queries
	googleClientId string
}

func NewServer() *http.Server {
	// Get config:
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Cannot parse environment variable `port` as int")
	}

	isDevelopment := false
	if os.Getenv("IS_DEVELOPMENT") != "" {
		isDevelopment, err = strconv.ParseBool(os.Getenv("IS_DEVELOPMENT"))
		if err != nil {
			panic("Cannot parse environment variable `is_development` as boolean")
		}
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to database: %s", err))
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientId == "" {
		panic("Cannot find environment variable `GOOGLE_CLIENT_ID`")
	}

	server := Server{
		config: Config{
			port:          port,
			isDevelopment: isDevelopment,
		},
		dbQueries:      database.New(db),
		googleClientId: googleClientId,
	}

	httpServer := &http.Server{
		Handler: server.RegisterRoutes(),
		Addr:    fmt.Sprintf(":%d", port),
	}

	return httpServer
}
