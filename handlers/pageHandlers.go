package handlers

import (
	"net/http"
	"outatime/util"
)

// Handle /
func Root(res http.ResponseWriter, req *http.Request) {

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

	err := util.RenderTemplate(res, "home", loggedIn, nil)
	if err != nil {
		util.LogError(err, "render")
	}
}
