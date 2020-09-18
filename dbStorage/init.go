package dbStorage

import (
	"database/sql"
	_ "github.com/lib/pq"

	"fmt"
	"github.com/GinkT/AvitoAssignmentJob/config"
	"log"
)

var db *sql.DB

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