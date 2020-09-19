package http_api

import (
	"encoding/json"
	"github.com/GinkT/AvitoAssignmentJob/db_storage"
	"net/http"
)

type balanceData struct {
	Id 			string 		`json:"id"`
	Balance 	float64		`json:"balance"`
	Currency 	string     	`json:"currency"`
}

func BalanceResponseInJson(w http.ResponseWriter, id string, balance float64, currency string) {
	if currency == "" {
		currency = "RUB"
	}
	dataToEncode := &balanceData{
		Id:			id,
		Balance:  	balance,
		Currency: 	currency,
	}
	json.NewEncoder(w).Encode(dataToEncode)
}

type balanceChangeData struct {
	Id 			string 		`json:"id"`
	Amount 		string     	`json:"amount"`
}

func BalanceChangeResponseInJson(w http.ResponseWriter, id , amount string) {
	dataToEncode := &balanceChangeData{
		Id:			id,
		Amount:		amount,
	}
	json.NewEncoder(w).Encode(dataToEncode)
}

type transferData struct {
	FromId 			string 		`json:"fromId"`
	ToId 			string 		`json:"toId"`
	Amount 			string     	`json:"amount"`
}

func TransferResponseInJson(w http.ResponseWriter, fromId , toId, amount string) {
	dataToEncode := &transferData{
		FromId:		fromId,
		ToId:		toId,
		Amount:		amount,
	}
	json.NewEncoder(w).Encode(dataToEncode)
}

func HistoryResponseInJson(w http.ResponseWriter, transactions []*db_storage.Transaction) {
	for _, transaction := range transactions {
		json.NewEncoder(w).Encode(transaction)
	}
}