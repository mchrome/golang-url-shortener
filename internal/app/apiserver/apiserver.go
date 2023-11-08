package apiserver

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	router *mux.Router
}

func New() *APIServer {
	return &APIServer{
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.router.HandleFunc("/", serveIndex)
	return http.ListenAndServe(":8000", s)
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "123")
}
