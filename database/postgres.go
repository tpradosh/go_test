package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"job_ping/models"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	/*connect to the db*/

	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "myuser")
	password := getEnv("DB_PASSWORD", "mypassword")
	dbname := getEnv("DB_NAME", "mydb")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitDB(ptr *sql.DB) {
	/*one time creation of each table in db*/

	tables := []string{
		"database/watches.sql",
		"database/alerts.sql",
		"database/results.sql",
	}

	for _, table := range tables {
		sqlCmd, err := ioutil.ReadFile(table)

		if err != nil {
			log.Fatal(err)
		}

		_, err = ptr.Exec(string(sqlCmd))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetAllWatches(db *sql.DB) ([]models.Watch, error) {
	/*gets all the watches from db*/

	rows, err := db.Query("SELECT id, url, interval_ms, expected_status, created_at FROM watches ORDER BY id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	watches := []models.Watch{}

	for rows.Next() {
		var w models.Watch
		if err := rows.Scan(&w.ID, &w.URL, &w.IntervalMS, &w.ExpectedStatus, &w.CreatedAt); err != nil {
			return nil, err
		}
		watches = append(watches, w)
	}

	return watches, nil
}
