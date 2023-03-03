package handlers

import (
	"fmt"
	"net/http"
	"outatime/util"
)

// Handle /
func Root(res http.ResponseWriter, req *http.Request) {

	// Stop logged in viewers from getting to main pages
	util.RedirectToUser(res, req)

	// Get user logged in status.
	loggedIn, _ := util.IsLoggedIn(req)

	err := util.RenderTemplate(res, "home", loggedIn, nil)
	if err != nil {
		util.LogError(err, "render")
	}
}

// Handle routes to /{{user_url}}/
func UserHome(res http.ResponseWriter, req *http.Request) {
	fmt.Println("at home")

	// Get user logged in status.
	// loggedIn, _ := util.IsLoggedIn(req)

}
