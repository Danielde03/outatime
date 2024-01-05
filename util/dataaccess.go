package util

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5432
	user         = "app"
	password     = "0u1Ot1m3"
	dbname       = "outatime"
	maxOpenConns = 10
	maxIdleCOnn  = 5
	maxLifetime  = 5 * time.Minute
)

var dbConn *sql.DB

// Get database connection
func GetConnection() *sql.DB {

	if dbConn != nil {
		return dbConn
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbConn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		dbConn = nil
		LogError(err, "database")
		os.Exit(1)
	}

	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleCOnn)
	dbConn.SetConnMaxLifetime(maxLifetime)

	return dbConn
}

// execute command to database
//
// return rows and error if one is thrown
//
// args can only be used for values, not fields
func DatabaseExecute(command string, args ...any) (*sql.Rows, error) {

	db := GetConnection()
	defer db.Close()

	// execute query and get results
	results, err := db.Query(command, args...)

	if err != nil {
		LogError(err, "database")
		return nil, err
	}

	return results, nil

}

// Check if a value is unique in a table
func IsUnique(value string, field string, table string) bool {

	rows, err := DatabaseExecute(fmt.Sprintf("SELECT \"%v\" FROM outatime.\"%v\" WHERE \"%v\" = $1", field, table, field), value)

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

// Return null string if string is empty
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
