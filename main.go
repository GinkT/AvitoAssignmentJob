package main

import (
	"database/sql"
	"fmt"
	"github.com/GinkT/AvitoAssignmentJob/config"
	"github.com/GinkT/AvitoAssignmentJob/db_storage"
	"log"
	"strconv"

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

	db, err := db_storage.NewDatabase(config)
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

	fmt.Println("Server is listening... HEEE YA YEE ")
	http.ListenAndServe(":8181", nil)

}


func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amount := r.URL.Query().Get("amount")

	result, err := db_storage.BalancePayment(Server.db, id, amount)
	if err != nil {
		log.Println("Error happened handling Balance Payment:", err)
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

	balance, err := db_storage.GetBalance(Server.db, id, "")
	if err != nil {
		log.Println("Error handling withdraw!", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Println("Error converting amount string to float!")
		w.WriteHeader(http.StatusConflict)
		return
	}

	var response string
	if balance < amountFloat {
		response = fmt.Sprintf("Request contains id: %s | amount: %s\nNot enough money!", id, amount)
	} else {
		result, err := db_storage.BalancePayment(Server.db, id, "-" + amount)
		if err != nil {
			log.Println("Error happened handling Balance Payment:", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		lastInsertId , _ := result.LastInsertId()
		response = fmt.Sprintf("Request contains id: %s | amount: %s\n[DB LOG] Rows affected: %d, Last Insert ID: %d",
			id, amount, rowsAffected, lastInsertId)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	fromId := r.URL.Query().Get("fromId")
	toId := r.URL.Query().Get("toId")
	amount := r.URL.Query().Get("amount")

	senderBalance, err := db_storage.GetBalance(Server.db, fromId, "")
	if err != nil {
		log.Println("Error handling transfer handler!", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Println("Error converting amount string to float!")
		w.WriteHeader(http.StatusConflict)
		return
	}

	var response string
	if senderBalance < amountFloat {
		response = fmt.Sprintf("Request contains fromId: %s | toId: %s | amount: %s\nNot enough money to transfer!", fromId, toId, amount)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, response)
		return
	}

	senderBalanceDecrease, err := db_storage.BalancePayment(Server.db, fromId, "-" + amount)
	if err != nil {
		log.Println("Error happened handling Transfer:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	senderRowsAffected, _ := senderBalanceDecrease.RowsAffected()
	senderLastInsertId , _ := senderBalanceDecrease.LastInsertId()

	receiverBalanceIncrease, err := db_storage.BalancePayment(Server.db, toId, amount)
	if err != nil {
		log.Println("Error happened handling Transfer:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	receiverRowsAffected, _ := receiverBalanceIncrease.RowsAffected()
	receiverLastInsertId , _ := receiverBalanceIncrease.LastInsertId()


	response = fmt.Sprintf("Request contains fromId: %s | toId: %s | amount: %s" +
		"\n[DB LOG] Sender Rows affected: %d, Sender Last Insert ID: %d" +
		"\nReceiver Rows affected: %d, Receiver Last Insert ID: %d",
		fromId, toId, amount, senderRowsAffected, senderLastInsertId, receiverRowsAffected, receiverLastInsertId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	currency := r.URL.Query().Get("currency")

	balance, err := db_storage.GetBalance(Server.db, id, currency)
	if err != nil {
		log.Println("Error handling Balance!", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	response := fmt.Sprintf("Request contains id: %s | currency: %s\nBalance in RUB: %f", id, currency, balance)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func haveEnoughMoney(db *sql.DB, id, amount string) {

}
