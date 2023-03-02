package handlers

import (
	"net/http"
	"outatime/util"
)

const auth_cookie_age = 60

// Log a user in
//
// get email and hashed password from form and match against database
//
// get token cookie and make sure it matches input token
//
// If good, set users auth_code in a cookie
func Login(res http.ResponseWriter, req *http.Request) {

	/* // TODO: get user's url based on ID and redirect. Make utility for this.
	loggedIn, id := util.IsLoggedIn(req)

	if loggedIn {
		http.Redirect(res, req, "/"+id+"/", http.StatusSeeOther)
	}
	*/

	// only handle post requests for login
	if req.Method == "GET" {
		http.Redirect(res, req, "/", http.StatusSeeOther)

	} else {

		// get form data
		email := req.FormValue("email")
		password := req.FormValue("password")
		token := req.FormValue("token")

		tokenCookie, err := req.Cookie("token")

		// if no token - sign of bad intent
		if err != nil {
			util.LogError(err, "cookies")
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}

		// remove token cookie
		http.SetCookie(res, &http.Cookie{
			Name:   "token",
			Path:   "/",
			MaxAge: -1,
		})

		// if token don't match - sign of bad intent
		if token != tokenCookie.Value {
			http.Redirect(res, req, "/", http.StatusSeeOther)
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

		// if no return val TODO: notify no account returned
		if auth_code == "" {
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}

		// store auth_code in cookie
		http.SetCookie(res, &http.Cookie{
			Name:   "auth_code",
			Value:  auth_code,
			Path:   "/",
			MaxAge: auth_cookie_age,
		})

		// TODO change redirect to user's home page - use util discussed in above
		http.Redirect(res, req, "/", http.StatusSeeOther)
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
