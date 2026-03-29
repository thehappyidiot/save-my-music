package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
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
	sessionStore   *sessions.CookieStore
}

func NewServer() *http.Server {
	// Get config:
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("Cannot parse environment variable `port` as int")
	}

	isDevelopment := false
	if os.Getenv("IS_DEVELOPMENT") != "" {
		isDevelopment, err = strconv.ParseBool(os.Getenv("IS_DEVELOPMENT"))
		if err != nil {
			panic("Cannot parse environment variable `is_development` as boolean")
		}
	}

	if isDevelopment {
		fmt.Print("Server is running in Development mode. Do NOT use in Production. Speak friend and enter: ")
		var confirmation string
		fmt.Scanln(&confirmation)
		if "mellon" != strings.ToLower(confirmation) {
			panic("You shall not pass 🧙")
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

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		panic("Cannot find environment variable `SESSION_KEY`")
	}
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		panic("Cannot find environment varaible `DOMAIN`")
	}

	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	}

	server := Server{
		config: Config{
			port:          port,
			isDevelopment: isDevelopment,
		},
		dbQueries:      database.New(db),
		googleClientId: googleClientId,
		sessionStore:   sessionStore,
	}

	httpServer := &http.Server{
		Handler: server.RegisterRoutes(),
		Addr:    fmt.Sprintf(":%d", port),
	}

	return httpServer
}
