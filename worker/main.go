package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"job_ping/database" //postgres.go filex
	"job_ping/models"
)

func checkSite(w models.Watch) models.Result {
	/*checks the status of the website and returns a result with info if down; nil if site is up*/

	start := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(w.URL)
	elapsed := time.Since(start)

	if err != nil { //returns this only on like an unexpected error with status 0 instead of like 200-400
		return models.Result{
			WatchID:        w.ID,
			Status:         0,
			ResponseTimeMS: int(elapsed.Milliseconds()),
			CreatedAt:      time.Now(),
		}
	}
	defer resp.Body.Close()

	return models.Result{
		WatchID:        w.ID,
		Status:         resp.StatusCode,
		ResponseTimeMS: int(elapsed.Milliseconds()),
		CreatedAt:      time.Now(),
	}
}

func startWatchLoop(w models.Watch) {
	ticker := time.NewTicker(60 * time.Second) //preset to check every min since interval inserition is bugged
	defer ticker.Stop()

	for {
		res := checkSite(w)
		fmt.Printf("site checking %s\n", w.URL)
		// Treat 2xx and 3xx as up; others (including 0) as down
		if res.Status < 200 || res.Status >= 400 {
			fmt.Printf("site is down %s\n", w.URL)
		}
		<-ticker.C
	}
}

func main() {

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// database.InitDB(db) //tables setup

	watches, err := database.GetAllWatches(db)
	if err != nil {
		log.Fatal("Failed to load watches:", err)
	}
	log.Printf("loaded %d watches", len(watches))

	for _, w := range watches {
		go startWatchLoop(w)
	}

	select {}

}
