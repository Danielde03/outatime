package util

import "net/http"

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
