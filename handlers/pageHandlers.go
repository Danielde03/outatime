package handlers

import (
	"net/http"
	"outatime/util"
)

// Handle /
func Root(res http.ResponseWriter, req *http.Request) {

	// Handle dynamic routes
	util.DynamicRoutes(res, req)

	// Get user logged in status.
	loggedIn, _ := util.IsLoggedIn(req)

	// open 404 page
	if req.URL.Path != "/" {
		err := util.RenderTemplate(res, "404", loggedIn, nil)
		if err != nil {
			util.LogError(err, "render")
		}
		return
	}

	// Stop logged in viewers from getting to main pages
	util.RedirectToUser(res, req)

	err := util.RenderTemplate(res, "home", loggedIn, nil)
	if err != nil {
		util.LogError(err, "render")
	}
}
