package main

import (
	dbStore "CorsGame/internal/store/sqlite"
	"log"
)

func main() {
	db, err := dbStore.OpenConnection()
	if err != nil {
		log.Fatal("could not connect to DB", err)
	}
	defer db.Close()

	err = dbStore.CreateTables(db)
	if err != nil {
		log.Fatal("could not create tables", err)
	}

}
