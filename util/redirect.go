package util

import "net/http"

// Redirect a logged out user to root page
//
// stops viewers from getting to user's pages
//
// If going to place allowed to go to, continue
func RedirectToRoot(res http.ResponseWriter, req *http.Request) {

	loggedIn, _ := IsLoggedIn(req)

	if !loggedIn {
		// redirect to home page
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

}

// Redirect a logged in user to his user-home page
//
// stops users from getting to main pages
//
// If going to place allowed to go to, continue
func RedirectToUser(res http.ResponseWriter, req *http.Request) {

	loggedIn, id := IsLoggedIn(req)

	if loggedIn {
		// redirect to user's home page
		http.Redirect(res, req, "/"+GetUserURL(id)+"/", http.StatusSeeOther)
		return
	}

}
