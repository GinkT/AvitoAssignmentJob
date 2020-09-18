package http_api

import (
	"database/sql"
	"fmt"
	"github.com/GinkT/AvitoAssignmentJob/db_storage"
	"log"
	"net/http"
	"strconv"
)

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amount := r.URL.Query().Get("amount")

	result, err := db_storage.BalanceChange(DB, id, amount)
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

	if result, err := haveEnoughMoney(DB, id, amount); !result {
		if err != nil {
			log.Println("Error checking money availible in Transfer handler:", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		response := fmt.Sprintf("Request contains fromId: %s | amount: %s\nNot enough money!", id, amount)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, response)
		return
	}

	result, err := db_storage.BalanceChange(DB, id, "-" + amount)
	if err != nil {
		log.Println("Error happened handling Balance Payment:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	lastInsertId , _ := result.LastInsertId()

	response := fmt.Sprintf("Request contains id: %s | amount: %s\n[DB LOG] Rows affected: %d, Last Insert ID: %d",
		id, amount, rowsAffected, lastInsertId)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	fromId := r.URL.Query().Get("fromId")
	toId := r.URL.Query().Get("toId")
	amount := r.URL.Query().Get("amount")

	if result, err := haveEnoughMoney(DB, fromId, amount); !result {
		if err != nil {
			log.Println("Error checking money availible in Transfer handler:", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		response := fmt.Sprintf("Request contains fromId: %s | toId: %s | amount: %s\nNot enough money!", fromId, toId, amount)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, response)
		return
	}

	senderBalanceDecrease, err := db_storage.BalanceChange(DB, fromId, "-" + amount)
	if err != nil {
		log.Println("Error happened handling Transfer:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	senderRowsAffected, _ := senderBalanceDecrease.RowsAffected()
	senderLastInsertId , _ := senderBalanceDecrease.LastInsertId()

	receiverBalanceIncrease, err := db_storage.BalanceChange(DB, toId, amount)
	if err != nil {
		log.Println("Error happened handling Transfer:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	receiverRowsAffected, _ := receiverBalanceIncrease.RowsAffected()
	receiverLastInsertId , _ := receiverBalanceIncrease.LastInsertId()

	response := fmt.Sprintf("Request contains fromId: %s | toId: %s | amount: %s" +
		"\n[DB LOG] Sender Rows affected: %d, Sender Last Insert ID: %d" +
		"\nReceiver Rows affected: %d, Receiver Last Insert ID: %d",
		fromId, toId, amount, senderRowsAffected, senderLastInsertId, receiverRowsAffected, receiverLastInsertId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	currency := r.URL.Query().Get("currency")

	balance, err := db_storage.GetBalance(DB, id, currency)
	if err != nil {
		log.Println("Error handling Balance!", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	var response string
	if currency != "" {
		response = fmt.Sprintf("Request contains id: %s | currency: %s\nBalance in %s: %f", id, currency, currency, balance)
	}
	response = fmt.Sprintf("Request contains id: %s \nBalance in RUB: %f", id, balance)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func haveEnoughMoney(db *sql.DB, id, amount string) (bool, error){
	balance, err := db_storage.GetBalance(db, id, "")
	if err != nil {
		return false, err
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Println("Error converting amount string to float!")
		return false, err
	}

	if balance < amountFloat {
		return false, nil
	}
	return true, nil
}