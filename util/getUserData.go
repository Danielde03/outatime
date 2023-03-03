package util

import (
	"errors"
	"net/http"
	"outatime/models"
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
		LogError(errors.New("GetUserURL() : database error"), "util")
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

// Get a user by ID
func GetUserById(id string) *models.User {

	rows, err := DatabaseExecute("SELECT user_name, user_avatar, \"isActive\", subscribers FROM outatime.user WHERE user_id = " + id + ";")

	if err != nil {
		LogError(err, "database")
		LogError(errors.New("GetUserById() : database error"), "util")
	}

	var user_name string
	var user_avatar string
	var user_active bool
	var user_subs int

	for rows.Next() {
		err := rows.Scan(&user_name, &user_avatar, &user_active, &user_subs)

		if err != nil {
			LogError(err, "database")
			LogError(errors.New("GetUserById() : database error"), "util")
		}
	}

	return &models.User{Name: user_name, Avatar: user_avatar, Active: user_active, Subs: user_subs}

}

// Get a user by URL
func GetUserByUrl(url string) *models.User {

	rows, err := DatabaseExecute("SELECT user_name, user_avatar, \"isActive\", subscribers FROM outatime.user WHERE user_url = '" + url + "';")

	if err != nil {
		LogError(err, "database")
		LogError(errors.New("GetUserByUrl() : database error"), "util")
	}

	var user_name string
	var user_avatar string
	var user_active bool
	var user_subs int

	for rows.Next() {
		err := rows.Scan(&user_name, &user_avatar, &user_active, &user_subs)

		if err != nil {
			LogError(err, "database")
			LogError(errors.New("GetUserByUrl() : database error"), "util")
		}
	}

	return &models.User{Name: user_name, Avatar: user_avatar, Active: user_active, Subs: user_subs}

}
