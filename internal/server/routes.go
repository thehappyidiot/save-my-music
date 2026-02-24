package server

import (
	"database/sql"
	"fmt"
	"html/template"
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

const INTERNAL_ERROR = "Something went wrong"

const LOGIN_SESSION_NAME = "user_login"
const AUTHENTICATED = "authenticated"
const USER_ID = "user_id"

// Set of endpoints that can be accessed without authentication
var OPEN_ENDPOINTS = map[string]bool{
	"/api/login":  true,
	"/api/health": true,
	"/app/login":  true,
}

func (server *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", server.getRoot)
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
			handler.ServeHTTP(w, req)
		}

		session, err := server.sessionStore.Get(req, LOGIN_SESSION_NAME)
		if err != nil {
			fmt.Printf("Error while parsing existing session: %s\n", err)
			session.IsNew = true // Treat it as a fresh session }
		}

		isAuthenticatedRaw := session.Values[AUTHENTICATED]
		if isAuthenticated, ok := isAuthenticatedRaw.(bool); !ok || !isAuthenticated || session.IsNew {
			// Redirect to login:
			http.Redirect(w, req, "/app/login", http.StatusSeeOther)
			return
		}

		// Authenticated user, continue
		handler.ServeHTTP(w, req)
	})
}

func (server *Server) getRoot(w http.ResponseWriter, req *http.Request) {
	homepageTemplate, err := template.ParseFiles("./frontend/index.html")
	if err != nil {
		http.Error(w, INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	// Get the cookie:
	session, err := server.sessionStore.Get(req, LOGIN_SESSION_NAME)
	if err != nil {
		http.Error(w, INTERNAL_ERROR, http.StatusInternalServerError)
		fmt.Printf("Error while fetching session store: %s\n", err)
		return
	}

	userIdRaw := session.Values[USER_ID]
	userId, ok := userIdRaw.(string)
	if !ok {
		http.Error(w, INTERNAL_ERROR, http.StatusInternalServerError)
		fmt.Print("Error while converting userId to String\n")
		return
	}
	err = homepageTemplate.Execute(w, userId)
	if err != nil {
		http.Error(w, INTERNAL_ERROR, http.StatusInternalServerError)
		fmt.Printf("Error while executing template: %s\n", err)
	}
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
	_, err = server.dbQueries.GetUserBySub(req.Context(), util.StringToNullString(payload.Subject))

	if err == sql.ErrNoRows {
		// New user
		server.dbQueries.UpsertUser(req.Context(), database.UpsertUserParams{
			GoogleSub:  util.StringToNullString(payload.Subject),
			Email:      util.InterfaceToNullString(payload.Claims["email"]),
			PictureUrl: payload.Claims["picture"],
			FullName:   util.InterfaceToNullString(payload.Claims["name"]),
			GivenName:  util.InterfaceToNullString(payload.Claims["given_name"]),
			FamilyName: util.InterfaceToNullString(payload.Claims["family_name"]),
		})
		return
	} else if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Login failed: %s", err),
			http.StatusInternalServerError,
		)
		fmt.Println("Error while querying user data: ", err)
	}

	// Set cookies:
	session, _ := server.sessionStore.Get(req, LOGIN_SESSION_NAME)
	session.Values[AUTHENTICATED] = true
	session.Values[USER_ID] = payload.Subject
	err = session.Save(req, w)
	if err != nil {
		http.Error(w, INTERNAL_ERROR, http.StatusInternalServerError)
		fmt.Printf("Error while saving cookies: %s\n", err.Error())
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
