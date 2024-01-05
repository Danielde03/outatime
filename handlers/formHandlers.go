package handlers

import (
	"errors"
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

		// check token
		if !util.CheckToken(res, req) {
			return
		}

		rows, err := util.DatabaseExecute("SELECT user_id FROM outatime.user WHERE user_email = $1 AND user_password = $2;", email, password)

		if err != nil {
			util.LogError(err, "database")
		}
		defer rows.Close()

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

	// check token
	if !util.CheckToken(res, req) {
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
	_, err = util.DatabaseExecute("INSERT INTO outatime.\"user\"(user_name, user_url, user_email, user_password, user_avatar, validate_code, \"isValid\", \"isActive\") VALUES ($1, $2, $3, $4, 'avatar.png', $5, true, true)", username, url, email, password, validation_code)
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

	// check token
	if !util.CheckToken(res, req) {
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

		// delete old image file if one exists
		if strings.TrimSpace(util.GetUserPage(id).Banner) != "" {

			err = os.Remove("./templates/public/images/" + util.GetUserById(id).Url + "/" + util.GetUserPage(id).Banner)
		}
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

// Update a user's account based on form input
func UpdateAccount(res http.ResponseWriter, req *http.Request) {

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

	// check token
	if !util.CheckToken(res, req) {
		return
	}

	// Get form values
	username := req.FormValue("user-name")
	url := req.FormValue("user-url")

	//TODO: validate data

	// return error if url has changed and is not unique.
	if url != util.GetUserURL(id) && !util.IsUnique(url, "user_url", "user") {
		http.Error(res, "URL is taken", 500)
		return
	}

	// get image for avatar
	req.ParseMultipartForm(32 << 20)
	imageFile, handler, err := req.FormFile("avatar")

	// don't worry about "no such file" errors
	if err != nil && err.Error() != "http: no such file" {
		util.LogError(err, "files")
		http.Error(res, "Image error", 500)
		return
	}

	avatar := ""

	// if no file added
	if err != nil && err.Error() == "http: no such file" {

		// set image directory to updated URL
		os.Rename("./templates/public/images/"+util.GetUserById(id).Url, "./templates/public/images/"+url)

		// update URL in image string
		newAvatarString := strings.Replace(util.GetUserById(id).Avatar, util.GetUserById(id).Url, url, 1)

		// update user without new avatar
		_, err = util.DatabaseExecute("UPDATE outatime.\"user\" SET user_name=$1, user_url=$2, user_avatar=$3 WHERE user_id=$4;", username, url, newAvatarString, id)
		if err != nil {
			util.LogError(err, "database")
			util.LogError(errors.New("URL in files and DB is out of sync. File: "+url+". DB: "+util.GetUserById(id).Url), "database")
			util.LogError(errors.New("URL in files and DB is out of sync. File: "+url+". DB: "+util.GetUserById(id).Url), "files")
			http.Error(res, "Database error", 500)
			panic("URL in files and DB is out of sync. File: " + url + ". DB: " + util.GetUserById(id).Url) // if URL and DB out of sync, end routine
		}

	} else { // if file is added

		// delete old image file, only if not the default avatar.png
		if util.GetUserById(id).Avatar != "avatar.png" {

			err = os.Remove("./templates/public/images/" + util.GetUserById(id).Avatar)
			if err != nil {
				util.LogError(err, "files")
				http.Error(res, "Image file deletion error", 500)
				return
			}

		}

		// set image directory name to updated URL
		os.Rename("./templates/public/images/"+util.GetUserById(id).Url, "./templates/public/images/"+url)

		// make image file based on new URL
		defer imageFile.Close()
		avatar = handler.Filename
		file, err := os.OpenFile("./templates/public/images/"+url+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			util.LogError(err, "files")
			http.Error(res, "Image file creation error", 500)
			return
		}
		defer file.Close()

		// put input image content into new image file
		io.Copy(file, imageFile)

		// update user with new avatar
		_, err = util.DatabaseExecute("UPDATE outatime.\"user\" SET user_name=$1, user_url=$2, user_avatar=$3  WHERE user_id=$4;", username, url, url+"/"+avatar, id)
		if err != nil {
			util.LogError(err, "database")
			util.LogError(errors.New("URL in files and DB is out of sync. File: "+url+". DB: "+util.GetUserById(id).Url), "database")
			util.LogError(errors.New("URL in files and DB is out of sync. File: "+url+". DB: "+util.GetUserById(id).Url), "files")
			http.Error(res, "Database error", 500)
			panic("URL in files and DB is out of sync. File: " + url + ". DB: " + util.GetUserById(id).Url) // if URL and DB out of sync, end routine
		}

	}

}

// Update a users's password. Check both fields match
func UpdatePassword(res http.ResponseWriter, req *http.Request) {

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

	// check token
	if !util.CheckToken(res, req) {
		return
	}

	// get form data
	password := req.FormValue("password")
	repeat := req.FormValue("rep-password")

	if password != repeat {
		http.Error(res, "Password values don't match", 500)
		return
	}

	// update user without new avatar
	_, err := util.DatabaseExecute("UPDATE outatime.\"user\" SET user_password=$1 WHERE user_id=$2;", password, id)
	if err != nil {
		util.LogError(err, "database")
		http.Error(res, "Database error", 500)
		return
	}

}

// Create a new event based on form values
func CreateEvent(res http.ResponseWriter, req *http.Request) {

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

	// check token
	if !util.CheckToken(res, req) {
		return
	}

	// get form data
	req.ParseMultipartForm(32 << 20)
	eventName := req.FormValue("name")
	eventStart := req.FormValue("start")
	eventEnd := req.FormValue("end")
	eventLocation := req.FormValue("location")
	imageFile, imageHandler, err := req.FormFile("image")
	// _, imageHandler, err := req.FormFile("image")
	eventDescr := req.FormValue("descr")
	eventTldr := req.FormValue("tldr")
	eventView := req.FormValue("view")

	// if user gives a file, and there is still an error
	if err != nil && err.Error() != "http: no such file" {
		util.LogError(err, "files")
		http.Error(res, "Image file error", 500)
		return
	}

	// if hidden, give a code
	eventCode := ""
	if eventView == "hidden" {

		// if eventCode is blank or is not unique, get a new code
		for eventCode == "" || !util.IsUnique(eventCode, "event_code", "event") {

			eventCode, err = util.RandomString(5, 5)

			if err != nil {
				util.LogError(err, "util")
				http.Error(res, "Code generation error", 500)
				return
			}

		}

	}

	// if there is an image
	if imageHandler != nil {

		// make file
		defer imageFile.Close()
		imageName := imageHandler.Filename
		file, err := os.OpenFile("./templates/public/images/"+util.GetUserById(id).Url+"/"+imageHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			util.LogError(err, "files")
			http.Error(res, "Image file creation error", 500)
			return
		}
		defer file.Close()

		// put image into new file
		io.Copy(file, imageFile)

		// save event
		_, err = util.DatabaseExecute("INSERT INTO outatime.event(user_id, event_name, \"isPublic\", event_tldr, event_descr, event_start, event_end, event_location, event_img, event_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);", id, eventName, eventView == "public", eventTldr, eventDescr, eventStart, eventEnd, eventLocation, imageName, util.NewNullString(eventCode))
		if err != nil {
			util.LogError(err, "database")
			http.Error(res, "Database error", 500)
			return
		}

		// confirm
		http.Error(res, "Event created", 200)

	} else {
		// save event
		_, err = util.DatabaseExecute("INSERT INTO outatime.event(user_id, event_name, \"isPublic\", event_tldr, event_descr, event_start, event_end, event_location, event_img, event_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);", id, eventName, eventView == "public", eventTldr, eventDescr, eventStart, eventEnd, eventLocation, "", util.NewNullString(eventCode))
		if err != nil {
			util.LogError(err, "database")
			http.Error(res, "Database error", 500)
			return
		}

		// confirm
		http.Error(res, "Event created", 200)

	}

}
