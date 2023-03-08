package handlers

import (
	"net/http"
	"os"
	"outatime/util"
	"strings"
)

// time a user will stay logged in (seconds)
const auth_cookie_age = 60 * 60 * 24

// Log a user in
//
// get email and hashed password from form and match against database
//
// get token cookie and make sure it matches input token
//
// If good, set users auth_code in a cookie
func Login(res http.ResponseWriter, req *http.Request) {

	// only handle post requests for login
	if req.Method == "GET" {
		http.Redirect(res, req, "/", http.StatusSeeOther)

	} else {

		// get form data
		email := req.FormValue("email")
		password := req.FormValue("password")
		token := req.FormValue("token")

		tokenCookie, err := req.Cookie("token")

		// remove token cookie
		http.SetCookie(res, &http.Cookie{
			Name:   "token",
			Path:   "/",
			MaxAge: -1,
		})

		// if no token - sign of bad intent
		if err != nil {
			util.LogError(err, "cookies")
			http.Error(res, "Token cookie not found", 500)
			return
		}

		// if token don't match - sign of bad intent
		if token != tokenCookie.Value {
			http.Error(res, "Token does not match", 500)
			return
		}

		rows, err := util.DatabaseExecute("SELECT user_id, auth_code FROM outatime.user WHERE user_email = '" + email + "' AND user_password = '" + password + "';")

		if err != nil {
			util.LogError(err, "database")
		}

		userId := ""
		auth_code := ""

		for rows.Next() {
			err := rows.Scan(&userId, &auth_code)

			if err != nil {
				util.LogError(err, "database")
			}
		}

		// if no return val
		if auth_code == "" {
			http.Error(res, "Email or password is invalid", 500)
			return
		}

		// store auth_code in cookie
		http.SetCookie(res, &http.Cookie{
			Name:   "auth_code",
			Value:  auth_code,
			Path:   "/",
			MaxAge: auth_cookie_age,
		})

		// redirect to user's home page
		http.Error(res, "Logged in", 200)
		return
	}

}

// Log the user out
func Logout(res http.ResponseWriter, req *http.Request) {

	// delete auth cookie
	http.SetCookie(res, &http.Cookie{
		Name:   "auth_code",
		Path:   "/",
		MaxAge: -1,
	})

	req.Header.Del("user_id")

	// redirect to home page
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// Create a new account
//
// Will assign given email, username and password,
// and will make codes and url.
//
// If good - login and return 200.
//
// If bad - return 500 and message.
func SignUp(res http.ResponseWriter, req *http.Request) {

	// get values from form
	// email := req.FormValue("email")
	// password := req.FormValue("password")
	username := req.FormValue("username")
	token := req.FormValue("token")

	tokenCookie, err := req.Cookie("token")

	// delete token cookie
	http.SetCookie(res, &http.Cookie{
		Name:   "token",
		Path:   "/",
		MaxAge: -1,
	})

	// if no token - sign of bad intent
	if err != nil {
		util.LogError(err, "cookies")
		http.Error(res, "Token cookie not found", 500)
		return
	}

	// if token don't match - sign of bad intent
	if token != tokenCookie.Value {
		http.Error(res, "Token does not match", 500)
		return
	}

	// ensure email is unique

	// make url - ensure is unique TODO: add more replace functions based on validation
	url := username
	url = strings.ReplaceAll(url, " ", "")
	url = strings.ReplaceAll(url, "?", "")

	// make auth_code and validate_code - if non unique error is returned, remake codes and reinsert account
	auth_code, err := util.RandomString(5, 5)
	if err != nil {
		util.LogError(err, "util")
		http.Error(res, "Error making auth_code", 500)
		return
	}

	validation_code, err := util.RandomString(5, 5)
	if err != nil {
		util.LogError(err, "util")
		http.Error(res, "Error making validation_code", 500)
		return
	}

	// make file in images. TODO: make happen after account is verified.
	err = os.Mkdir("./templates/public/images/"+url, 0655)
	if err != nil {
		util.LogError(err, "files")
		http.Error(res, "Error making image file", 500)
		return
	}

	http.Error(res, "Made it to the end!", 500)

}
