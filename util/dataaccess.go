package util

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "@dm1n"
	dbname   = "outatime"
)

// execute command to database
//
// Return rows and error if one is thrown
func DatabaseExecute(command string) (*sql.Rows, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	// execute query and get results
	results, err := db.Query(command)

	if err != nil {
		return nil, err
	}

	// defer moves it to end of function - will happen last
	defer db.Close()

	return results, nil

}
