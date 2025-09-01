package main

import (
    "fmt"
    //"log"
    //"net/http"
    "time"
    _ "github.com/lib/pq"

	"job_ping/database" //postgres.go file
)

//urls to monitor
type Watch struct {
	ID int
	URL string
	Interval int
	ExpectedStatus int
}

//Resutl from checking each 'Watch'ed url
type Result struct {
	WatchID int
	Status int
	LatencyMS int64
	Success bool
	CheckedAt time.Time
}

func main() {

	db, _ := database.ConnectDB()
	database.InitDB(db) //tables setup
	fmt.Println("Hello, World!")


}