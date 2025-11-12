package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"google.golang.org/api/idtoken"
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

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idTokenString := req.FormValue("credential")

	payload, err := idtoken.Validate(req.Context(), idTokenString, server.googleClientId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(TYPE, TYPE_PLAIN)
	io.WriteString(w, "Welcome "+fmt.Sprintf("%v", payload.Claims["name"]))
}
