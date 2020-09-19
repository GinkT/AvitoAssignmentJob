package http_api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

var DB *sql.DB

var rout *mux.Router

// Создание роутера и маунт ссылок
func NewRouter(db *sql.DB) *mux.Router {
	if rout != nil {
		return rout
	}

	DB = db

	router := mux.NewRouter()
	router.HandleFunc("/payment", PaymentHandler)
	router.HandleFunc("/withdraw", WithdrawHandler)
	router.HandleFunc("/transfer", TransferHandler)
	router.HandleFunc("/balance", BalanceHandler)
	router.HandleFunc("/history", HistoryHandler)
	http.Handle("/",router)
	return router
}

