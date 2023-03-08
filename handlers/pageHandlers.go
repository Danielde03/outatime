package handlers

import (
	"net/http"
	"outatime/models"
	"outatime/util"
)

// TODO: send seperate data to nav bar as to body

// Handle /
func Root(res http.ResponseWriter, req *http.Request) {

	// Stop logged in viewers from getting to main pages
	util.RedirectToUser(res, req)

	// Get user logged in status.
	loggedIn, id := util.IsLoggedIn(req)
	var data models.PageData

	if loggedIn {
		data.NavUser = *util.GetUserById(id)
	}

	err := util.RenderTemplate(res, "home", loggedIn, data)
	if err != nil {
		util.LogError(err, "render")
	}
}

// Handle routes to /{{user_url}}/
func UserHome(res http.ResponseWriter, req *http.Request, user_url string) {

	// Get user logged in status.
	loggedIn, id := util.IsLoggedIn(req)
	var data models.PageData

	if loggedIn {
		data.NavUser = *util.GetUserById(id)
	}

	data.PageUser = *util.GetUserByUrl(user_url)
	data.UserPage = *util.GetUserPage(util.GetUserId(user_url))

	// TODO: make page viewed if not owner, make page editable if owner

	err := util.RenderTemplate(res, "user-home", loggedIn, data)
	if err != nil {
		util.LogError(err, "render")
	}

}
