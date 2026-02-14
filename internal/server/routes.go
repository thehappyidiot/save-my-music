package server

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/thehappyidiot/save-my-music/internal/database"
	"github.com/thehappyidiot/save-my-music/internal/util"
)

const TYPE = "Content-Type"
const TYPE_HTML = "text/html; charset=utf-8"
const TYPE_PLAIN = "text/plain; charset=utf-8"
const TYPE_JSON = "text/json; charset=utf-8"

// Set of endpoints that can be accessed without authentication
var OPEN_ENDPOINTS = map[string]bool{
	"/api/login":  true,
	"/api/health": true,
	"/app/login":  true,
}

func (server *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./frontend")))
	mux.HandleFunc("GET /app/login", server.getLogin)
	mux.HandleFunc("GET /api/health", server.getHealth)
	mux.HandleFunc("POST /api/login", server.postLogin)

	if server.config.isDevelopment {
		fmt.Print("Server is running in Development mode. Do NOT use in Production. Speak friend and enter: ")
		var confirmation string
		fmt.Scanln(&confirmation)
		if "mellon" != strings.ToLower(confirmation) {
			panic("You shall not pass 🧙")
		}
	}
	return server.middlewareHandler(mux)
}

func (server *Server) middlewareHandler(handler http.Handler) http.Handler {
	return server.middlewareLogger(server.middlewareAuth(handler))
}

func (server *Server) middlewareLogger(handler http.Handler) http.Handler {
	if !server.config.isDevelopment {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(res))
		handler.ServeHTTP(w, req)
	})
}

func (server *Server) middlewareAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, ok := OPEN_ENDPOINTS[req.URL.EscapedPath()]; ok {
			fmt.Println("Not checking auth, much wow")
			handler.ServeHTTP(w, req)
		}
		fmt.Println("Checking auth, such secure")
		handler.ServeHTTP(w, req)
	})
}

func (server *Server) getLogin(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./frontend/login.html")
}

func (server *Server) getHealth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(TYPE, TYPE_PLAIN)
	io.WriteString(w, "OK")
}

func (server *Server) postLogin(w http.ResponseWriter, req *http.Request) {
	payload, err := util.ValidateGoogleAuthRequest(req, server.googleClientId)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Login failed: %s", err),
			http.StatusUnauthorized,
		)
		return
	}

	// Get user
	user, err := server.dbQueries.GetUserBySub(req.Context(), util.StringToNullString(payload.Subject))

	var message string

	if err == sql.ErrNoRows {
		// New user
		message = "Welcome to hotel transylvania " + fmt.Sprintf("%v", payload.Claims["name"])
		server.dbQueries.UpsertUser(req.Context(), database.UpsertUserParams{
			GoogleSub:  util.StringToNullString(payload.Subject),
			Email:      util.InterfaceToNullString(payload.Claims["email"]),
			PictureUrl: payload.Claims["picture"],
			FullName:   util.InterfaceToNullString(payload.Claims["name"]),
			GivenName:  util.InterfaceToNullString(payload.Claims["given_name"]),
			FamilyName: util.InterfaceToNullString(payload.Claims["family_name"]),
		})
	} else if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Login failed: %s", err),
			http.StatusInternalServerError,
		)
		fmt.Println("Error while querying user data: ", err)
	} else {
		// User found:
		if user.GivenName.Valid {
			message = "Good to see you Mr. Wick, " + user.GivenName.String
		} else {
			message = "Bengan AF"
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(TYPE, TYPE_PLAIN)
	io.WriteString(w, message)
}
