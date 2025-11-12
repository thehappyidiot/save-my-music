package server

import (
	"io"
	"net/http"
)

const TYPE = "Content-Type"
const TYPE_HTML = "text/html; charset=utf-8"
const TYPE_PLAIN = "text/plain; charset=utf-8"
const TYPE_JSON = "text/json; charset=utf-8"

func (server *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", server.getHealth)

	return mux
}

func (server *Server) getHealth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(TYPE, TYPE_PLAIN)
	io.WriteString(w, "OK")
}
