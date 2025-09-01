package database

import (
    "log"
    _ "github.com/lib/pq"
    "io/ioutil"
	"database/sql"

)


func ConnectDB() (*sql.DB, error) {
	/*connect to the db*/

	connStr := "host=localhost port=5431 user=myuser password=mypassword dbname = mydb sslmode = disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB(ptr *sql.DB) {
	/*one time creation of each table in db*/

	tables := []string {
		"database/watches.sql",
		"database/alerts.sql",
		"database/results.sql",
	}

	for _, table := range tables{
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