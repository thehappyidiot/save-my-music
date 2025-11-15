package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set(TYPE, TYPE_PLAIN)
	io.WriteString(w, "Welcome "+fmt.Sprintf("%v", payload.Claims["name"]))
}
