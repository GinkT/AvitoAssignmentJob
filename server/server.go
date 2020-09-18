package server

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	db *sql.DB
	router *mux.Router
}

func NewServer(db *sql.DB, router *mux.Router) *Server {
	return &Server{
		db:db,
		router:router,
	}
}

func (server *Server) Run() {
	log.Println("Starting to listen and serve at :8181")
	http.ListenAndServe(":8181", nil)
}