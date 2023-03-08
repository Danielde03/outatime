package util

import (
	"database/sql"
	"errors"
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

// TODO: get rid of concatination for queries
// TODO: set up new account on DB for the app to use, not "postgres @dm1n"

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

// Check if a value is unique in a table
func IsUnique(value string, field string, table string) bool {

	rows, err := DatabaseExecute("SELECT " + field + " FROM outatime." + table + " WHERE " + field + " = '" + value + "'")

	if err != nil {
		LogError(err, "database")
		LogError(errors.New("IsUnique() : database error"), "util")
	}

	// get number of rows returned
	count := 0
	for rows.Next() {
		count += 1
	}

	// If rows are returned, it is in database, return false
	return count == 0
}
