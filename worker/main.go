package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	_ "github.com/lib/pq"

	"job_ping/database" //postgres.go filex
	"job_ping/models"
)




/*

double check this part 


func sendEmail(r models.Result, url string) {
	EMAIL := "testetstdsa@tset@gmail.com"

	// SMTP configuration (using Gmail as example)
	from := "your-email@gmail.com"
	password := "your-app-password" // Use app password, not regular password
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Email content
	subject := "Site Down Alert"
	body := fmt.Sprintf("Site is down!\n\nURL: %s\nStatus: %d\nResponse Time: %d ms\nTime: %s",
		url, r.Status, r.ResponseTimeMS, r.CreatedAt.Format("2006-01-02 15:04:05"))

	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", EMAIL, subject, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{EMAIL}, []byte(message))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	} else {
		log.Printf("Email sent successfully for down site: %s", url)
	}
}
*/
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
		// Treat 2xx and 3xx as up else down
		if res.Status < 200 || res.Status >= 400 {
			fmt.Printf("site is down %s\n", w.URL)

			//send to email
			//sendEmail(res, w.URL)
			//store in db
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
