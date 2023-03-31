package handlers

import (
	"net/http"
	"outatime/models"
	"outatime/util"
	"strings"
)

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
	loggedIn, loggedIn_id := util.IsLoggedIn(req)
	var data models.PageData

	data.PageUser = *util.GetUserByUrl(user_url)
	data.UserPage = *util.GetUserPage(util.GetUserId(user_url))

	if loggedIn {

		// if logged in user's id does not match the page's user id, redirect
		if loggedIn_id != util.GetUserId(user_url) {
			util.RedirectToUser(res, req)
		}

		// if logged in to own, make own page editable
		data.NavUser = *util.GetUserById(loggedIn_id)
		err := util.RenderTemplate(res, "user-home-edit", loggedIn, data)
		if err != nil {
			util.LogError(err, "render")
		}
	} else {

		// If viewers are accessing page

		data.UserPage.About = strings.Replace(data.UserPage.About, "\r\n", "<br>", -1)
		err := util.RenderTemplate(res, "user-home", loggedIn, data)
		if err != nil {
			util.LogError(err, "render")
		}

	}

}

// Handle /{{user_url}}/live-view
func UserHomeLiveView(res http.ResponseWriter, req *http.Request, user_url string) {

	// Get user logged in status.
	loggedIn, loggedIn_id := util.IsLoggedIn(req)
	var data models.PageData

	// if not logged in, or not the owner, go to root
	if !loggedIn || (loggedIn && loggedIn_id != util.GetUserId(user_url)) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	data.PageUser = *util.GetUserByUrl(user_url)
	data.NavUser = *util.GetUserById(loggedIn_id)
	data.UserPage = *util.GetUserPage(util.GetUserId(user_url))

	data.UserPage.About = strings.Replace(data.UserPage.About, "\r\n", "<br>", -1)

	util.RenderTemplate(res, "user-home", loggedIn, data)

}
