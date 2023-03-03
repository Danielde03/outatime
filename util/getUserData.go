package util

import (
	"errors"
	"net/http"
)

// Determine if a user is logged in.
//
// Returns true if so and a sting of the id.
func IsLoggedIn(req *http.Request) (bool, string) {

	id := req.Header.Get("user_id")

	if len(id) > 0 {
		return true, id
	}

	return false, id

}

// Get a user's URL based on the ID
//
// Empty URL means no user at that ID
func GetUserURL(id string) string {

	rows, err := DatabaseExecute("SELECT user_url FROM outatime.user WHERE user_id = " + id + ";")

	if err != nil {
		LogError(err, "database")
	}

	userURL := ""

	for rows.Next() {
		err := rows.Scan(&userURL)

		if err != nil {
			LogError(err, "database")
			LogError(errors.New("GetUserURL() : database error"), "util")
		}
	}

	defer rows.Close()
	return userURL

}
