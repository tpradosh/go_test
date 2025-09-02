package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"

	"job_ping/database" //postgres.go file
	"job_ping/models"

	
)


func getWatches(db *sql.DB) ([]models.Watch, error) {

	return nil, nil
}


func main() {

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// database.InitDB(db) //tables setup

	//watches, err := database.GetAllWatches(db)
}
