package handlers

import (
	"io"
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
		return

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
		if token != tokenCookie.Value || strings.TrimSpace(token) == "" {
			http.Error(res, "Token does not match", 500)
			return
		}

		rows, err := util.DatabaseExecute("SELECT user_id FROM outatime.user WHERE user_email = $1 AND user_password = $2;", email, password)

		if err != nil {
			util.LogError(err, "database")
		}

		userId := ""
		auth_code := ""

		for rows.Next() {
			err := rows.Scan(&userId)

			if err != nil {
				util.LogError(err, "database")
			}
		}

		// if no return val
		if userId == "" {
			http.Error(res, "Email or password is invalid", 500)
			return
		}

		// make auth_code and validate_code - if non unique, get new auth_code
		auth_code, err = util.RandomString(5, 5)

		for !util.IsUnique(auth_code, "auth_code", "user") {
			auth_code, err = util.RandomString(5, 5)
		}

		if err != nil {
			util.LogError(err, "util")
			return
		}

		// store auth_code in database
		_, err = util.DatabaseExecute("UPDATE outatime.\"user\" SET auth_code = $1 WHERE user_id = $2", auth_code, userId)
		if err != nil {
			util.LogError(err, "database")
			http.Error(res, "Database error", 500)
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

	// only handle post requests for logout
	if req.Method == "GET" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	} else {

		// delete auth cookie
		http.SetCookie(res, &http.Cookie{
			Name:   "auth_code",
			Path:   "/",
			MaxAge: -1,
		})

		// nullify auth_code in database TODO: nullify when cookie dies
		_, err := util.DatabaseExecute("UPDATE outatime.\"user\" SET auth_code = NULL where user_id = $1", req.Header.Get("user_id"))
		if err != nil {
			util.LogError(err, "database")
			http.Error(res, "Database error", 500)
			return
		}

		// delete from header
		req.Header.Del("user_id")

		// redirect to home page
		http.Redirect(res, req, "/", http.StatusSeeOther)

	}

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

	if req.Method == "GET" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// get values from form
	email := req.FormValue("email")
	password := req.FormValue("password")
	username := req.FormValue("username")
	url := req.FormValue("url")
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

	// ensure email is unique - if not, return error
	if !util.IsUnique(email, "user_email", "user") {
		http.Error(res, "Email is already in use", 500)
		return
	}

	// ensure URL is unique - if not, return error
	if !util.IsUnique(url, "user_url", "user") {
		http.Error(res, "URL is already in use", 500)
		return
	}

	validation_code, err := util.RandomString(5, 5)
	if err != nil {
		util.LogError(err, "util")
		http.Error(res, "Error making validation_code", 500)
		return
	}

	// add to database TODO: take out isActive and isValid
	_, err = util.DatabaseExecute("INSERT INTO outatime.\"user\"(user_name, user_url, user_email, user_password, user_avatar, validate_code, \"isValid\", \"isActive\") VALUES ($1, $2, $3, $4, ' ', $5, true, true)", username, url, email, password, validation_code)
	if err != nil {
		util.LogError(err, "database")
		http.Error(res, "Database error", 500)
		return
	}

	// make file in images. TODO: make happen after account is verified.
	err = os.Mkdir("./templates/public/images/"+url, 0655)
	if err != nil {
		util.LogError(err, "files")
		http.Error(res, "Error making image file", 500)
		return
	}

	http.Error(res, "Account made", http.StatusCreated)

}

// Update a user's page based on form input
func UpdatePage(res http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// check user is logged in -extra level of security
	loggedIn, id := util.IsLoggedIn(req)

	if !loggedIn {
		http.Error(res, "No user logged in", 500)
		return
	}

	// get form data
	about := req.FormValue("about")
	public := req.FormValue("isPublic")

	// get token data
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
	if token != tokenCookie.Value || strings.TrimSpace(token) == "" {
		http.Error(res, "Token does not match", 500)
		return
	}

	// TODO: validate data

	// if image is new, add image to file

	req.ParseMultipartForm(32 << 20)
	imageFile, handler, err := req.FormFile("banner")

	// don't worry about "no such file" errors
	if err != nil && err.Error() != "http: no such file" {
		util.LogError(err, "files")
		http.Error(res, "Image error", 500)
		return
	}

	banner := ""

	// if no file added
	if err != nil && err.Error() == "http: no such file" {

		// update user without new banner
		_, err = util.DatabaseExecute("UPDATE outatime.user_page SET \"aboutUs\"=$1, \"isPublic\"=$2 WHERE user_id=$3;", about, public == "true", id)
		if err != nil {
			util.LogError(err, "database")
			http.Error(res, "Database error", 500)
			return
		}

	} else { // if file is added

		// make file
		defer imageFile.Close()
		banner = handler.Filename
		file, err := os.OpenFile("./templates/public/images/"+util.GetUserById(id).Url+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			util.LogError(err, "files")
			http.Error(res, "Image file creation error", 500)
			return
		}
		defer file.Close()

		// put image into new file
		io.Copy(file, imageFile)

		// delete old image file
		os.Remove("./templates/public/images/" + util.GetUserById(id).Url + "/" + util.GetUserPage(id).Banner)
		if err != nil {
			util.LogError(err, "files")
			http.Error(res, "Image file deletion error", 500)
			return
		}

		// update user with new banner
		_, err = util.DatabaseExecute("UPDATE outatime.user_page SET \"aboutUs\"=$1, banner=$2, \"isPublic\"=$3 WHERE user_id=$4;", about, banner, public == "true", id)
		if err != nil {
			util.LogError(err, "database")
			http.Error(res, "Database error", 500)
			return
		}

	}

}
