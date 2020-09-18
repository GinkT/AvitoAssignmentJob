package main

import (
	"database/sql"
	"github.com/GinkT/AvitoAssignmentJob/config"
	"github.com/GinkT/AvitoAssignmentJob/dbStorage"
	"log"
)

type server struct {
	db *sql.DB
}

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

	_ = newServer(db)

}