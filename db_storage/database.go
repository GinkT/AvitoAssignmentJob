package db_storage

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

// Подключение к БД
func NewDatabase(dbHost string, dbPort int, dbUser string, dbPassword string, dbBase string) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbBase)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	log.Println("Database was successfully connected!")
	return db, nil
}
