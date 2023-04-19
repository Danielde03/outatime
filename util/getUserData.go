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

	rows, err := DatabaseExecute("SELECT user_url FROM outatime.user WHERE user_id = $1;", id)

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

// Get a user's ID based on the URL
//
// Empty URL means no user at that ID
func GetUserId(url string) string {

	rows, err := DatabaseExecute("SELECT user_id FROM outatime.user WHERE user_url = $1;", url)

	if err != nil {
		LogError(err, "database")
		LogError(errors.New("GetUserId() : database error"), "util")
	}

	userId := ""

	for rows.Next() {
		err := rows.Scan(&userId)

		if err != nil {
			LogError(err, "database")
			LogError(errors.New("GetUserId() : database error"), "util")
		}
	}

	defer rows.Close()
	return userId

}

// Get a user by ID
func GetUserById(id string) *models.User {

	rows, err := DatabaseExecute("SELECT user_name, user_url, user_avatar, \"isActive\", subscribers, events FROM outatime.user WHERE user_id = $1;", id)

	if err != nil {
		LogError(err, "database")
		LogError(errors.New("GetUserById() : database error"), "util")
	}

	var user_name string
	var user_url string
	var user_avatar string
	var user_active bool
	var user_subs int
	var user_events int

	for rows.Next() {
		err := rows.Scan(&user_name, &user_url, &user_avatar, &user_active, &user_subs, &user_events)

		if err != nil {
			LogError(err, "database")
			LogError(errors.New("GetUserById() : database error"), "util")
		}
	}

	return &models.User{Name: user_name, Url: user_url, Avatar: user_avatar, Active: user_active, Subs: user_subs, Events: user_events}

}

// Get a user by URL
func GetUserByUrl(url string) *models.User {

	return GetUserById(GetUserId(url))

}

// Get the page data of a user based on user id
func GetUserPage(user_id string) *models.UserPage {

	rows, err := DatabaseExecute("SELECT \"aboutUs\", banner, \"isPublic\" FROM outatime.user_page WHERE user_id = $1;", user_id)

	if err != nil {
		LogError(err, "database")
		LogError(errors.New("GetUserPage() : database error"), "util")
	}

	var about string
	var banner string
	var public bool

	for rows.Next() {
		err := rows.Scan(&about, &banner, &public)

		if err != nil {
			LogError(err, "database")
			LogError(errors.New("GetUserPage() : database error"), "util")
		}
	}

	page := &models.UserPage{About: about, Banner: banner, Public: public}

	defer rows.Close()
	return page

}
