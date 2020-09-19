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

	log.Printf("[HTTP] Got a payment request, contains: id: %s | amount: %s\n", id, amount)

	result, err := db_storage.BalanceChange(DB, id, amount)
	if err != nil {
		log.Println("Error handling Payment!", err)
		fmt.Fprint(w, "Error handling Payment!", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("[DB] Payment handled. Rows Affected: %d\n", rowsAffected)

	db_storage.AddTransaction(DB, "payment", "", id, amount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	BalanceChangeResponseInJson(w, id, amount)
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amount := r.URL.Query().Get("amount")

	log.Printf("[HTTP] Got a withdraw request, contains: id: %s | amount: %s\n", id, amount)

	if result, err := haveEnoughMoney(DB, id, amount); !result {
		if err != nil {
			log.Println("Error checking money availible in Transfer handler:", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		fmt.Fprint(w, "Error handling Withdraw! Not enough money!")
		w.WriteHeader(http.StatusConflict)
		return
	}

	result, err := db_storage.BalanceChange(DB, id, "-" + amount)
	if err != nil {
		log.Println("Error happened handling Balance Payment:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("[DB] Payment handled. Rows Affected: %d\n", rowsAffected)

	db_storage.AddTransaction(DB, "withdraw", id, "", amount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	BalanceChangeResponseInJson(w, id, "-" + amount)
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	fromId := r.URL.Query().Get("fromId")
	toId := r.URL.Query().Get("toId")
	amount := r.URL.Query().Get("amount")

	log.Printf("[HTTP] Got a transfer request, contains: fromId: %s | toId: %s | amount: %s\n", fromId, toId, amount)

	if result, err := haveEnoughMoney(DB, fromId, amount); !result {
		if err != nil {
			log.Println("Error checking money availible in Transfer handler:", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		fmt.Fprint(w, "Error handling Withdraw! Not enough money!")
		w.WriteHeader(http.StatusConflict)
		return
	}

	senderBalanceDecrease, err := db_storage.BalanceChange(DB, fromId, "-" + amount)
	if err != nil {
		log.Println("Error happened handling Transfer:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	senderRowsAffected, _ := senderBalanceDecrease.RowsAffected()

	receiverBalanceIncrease, err := db_storage.BalanceChange(DB, toId, amount)
	if err != nil {
		log.Println("Error happened handling Transfer:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	receiverRowsAffected, _ := receiverBalanceIncrease.RowsAffected()

	log.Printf("[DB LOG] Sender Rows affected: %d | Receiver Rows affected: %d\n",
		senderRowsAffected, receiverRowsAffected)

	db_storage.AddTransaction(DB, "transfer", fromId, toId, amount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	TransferResponseInJson(w, fromId, toId, amount)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	currency := r.URL.Query().Get("currency")

	log.Printf("[HTTP] Got a balance request, contains: id: %s | currencyu: %s\n", id, currency)

	balance, err := db_storage.GetBalance(DB, id, currency)
	if err != nil {
		log.Println("Error handling Balance!", err)
		fmt.Fprint(w, "Error handling Balance!", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	BalanceResponseInJson(w, id, balance, currency)
}

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	sortBy := r.URL.Query().Get("sortBy")
	orderBy := r.URL.Query().Get("orderBy")

	log.Printf("[HTTP] Got a history request, contains: id: %s | sortBy: %s | orderBy: %s\n", id, sortBy, orderBy)

	transactions := db_storage.GetHistoryForId(DB, id, sortBy, orderBy)
	w.WriteHeader(http.StatusOK)
	HistoryResponseInJson(w, transactions)
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