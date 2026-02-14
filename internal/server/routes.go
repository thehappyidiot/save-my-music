package server

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/thehappyidiot/save-my-music/internal/database"
	"github.com/thehappyidiot/save-my-music/internal/util"
)

const TYPE = "Content-Type"
const TYPE_HTML = "text/html; charset=utf-8"
const TYPE_PLAIN = "text/plain; charset=utf-8"
const TYPE_JSON = "text/json; charset=utf-8"

func (server *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./frontend")))
	mux.HandleFunc("GET /api/health", server.getHealth)
	mux.HandleFunc("POST /api/login", server.postLogin)

	if !server.config.isDevelopment {
		return mux
	}

	fmt.Print("Server is running in Development mode. DO NOT use this in Production! Confirm by typing 'YOLO': ")
	var confirmation string
	fmt.Scanln(&confirmation)
	if "YOLO" != confirmation {
		panic("You should not be running in Development mode, change the config")
	}

	return server.middlewareLogger(mux)
}

func (server *Server) middlewareLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(res))
		handler.ServeHTTP(w, req)
	})
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
