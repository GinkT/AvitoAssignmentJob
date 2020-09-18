package main

import (
	"database/sql"
	"fmt"
	"github.com/GinkT/AvitoAssignmentJob/config"
	"github.com/GinkT/AvitoAssignmentJob/dbStorage"
	"log"

	"net/http"
	"github.com/gorilla/mux"
)

type server struct {
	db *sql.DB
}

var Server *server

func newServer(db *sql.DB) *server {
	return &server{
		db:db,
	}
}

func main () {
	config := &config.ConfigStruct{
		"db",
		5432,
		"postgres",
		"qwerty",
		"UserBilling",
	}

	db, err := dbStorage.NewDatabase(config)
	if err != nil {
		log.Fatalln("Database was not connected!")
	}

	Server = newServer(db)

	router := mux.NewRouter()
	router.HandleFunc("/payment", PaymentHandler)
	router.HandleFunc("/withdraw", WithdrawHandler)
	router.HandleFunc("/transfer", TransferHandler)
	router.HandleFunc("/balance", BalanceHandler)
	http.Handle("/",router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}


func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amount := r.URL.Query().Get("amount")

	result, err := dbStorage.BalancePayment(Server.db, id, amount)
	if err != nil {
		log.Println("Error happened hadling Balance Payment:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	lastInsertId , _ := result.LastInsertId()
	response := fmt.Sprintf("Request contains id: %s | amount: %s\n[DB LOG] Rows affected: %d, Last Insert ID: %d", id, amount, rowsAffected, lastInsertId)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amount := r.URL.Query().Get("amount")
	response := fmt.Sprintf("Request contains id: %s | amount: %s", id, amount)
	fmt.Fprint(w, response)
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amount := r.URL.Query().Get("amount")
	response := fmt.Sprintf("Request contains id: %s | amount: %s", id, amount)
	fmt.Fprint(w, response)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	currency := r.URL.Query().Get("currency")
	response := fmt.Sprintf("Request contains id: %s | currency: %s", id, currency)
	fmt.Fprint(w, response)
}

