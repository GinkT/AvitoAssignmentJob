package dbStorage

import (
	"database/sql"
	_ "github.com/lib/pq"

	"fmt"
	"github.com/GinkT/AvitoAssignmentJob/config"
	"log"
)

var db *sql.DB

type servDB sql.DB

// Подключение к БД
func NewDatabase(config *config.ConfigStruct) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	log.Println("Database was successfully connected!")
	return db, nil
}

/* Пополнение баланса. В случае если юзера не было в таблице - создаётся новая запись.
   Если юзер был - выполняет UPDATE баланса и прибавляет сумму пополнения*/
func BalancePayment(db *sql.DB, id, amount string) (sql.Result, error) {
	sqlStatement := `
			INSERT INTO public."users"
			VALUES($1, $2)
			ON CONFLICT ("id")
			DO
			UPDATE SET "balance" = "users"."amount" + $2
		`

	res, err := db.Exec(sqlStatement, id, amount)
	if err != nil {
		return nil, err
	}

	return res, nil
}