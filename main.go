package main

import (
	"github.com/GinkT/AvitoAssignmentJob/db_storage"
	"github.com/GinkT/AvitoAssignmentJob/http_api"
	"github.com/GinkT/AvitoAssignmentJob/server"
	"log"
)

const (
	dbHost = "db"
	dbPort = 5432
	dbUser = "postgres"
	dbPassword = "qwerty"
	dbBase = "UserBilling"
)


func main () {
	db, err := db_storage.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbBase)
	if err != nil {
		log.Fatalln("Database was not connected!")
	}

	router := http_api.NewRouter(db)

	server := server.NewServer(db, router)
	server.Run()
}


