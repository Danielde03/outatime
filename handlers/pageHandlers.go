package handlers

// General page handlers

import (
	"net/http"
	"outatime/models"
	"outatime/util"
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

// Handle /hosts
func Hosts(res http.ResponseWriter, req *http.Request) {

	// Stop logged in viewers from getting to main pages
	util.RedirectToUser(res, req)

	// get all public and active users

	var data models.PageData
	var name string
	var url string
	var avatar string
	var subs int
	var events int

	// get data
	rows, err := util.DatabaseExecute("SELECT user_name, user_url, user_avatar, subscribers, events FROM outatime.\"user\" JOIN outatime.\"user_page\" ON \"user\".user_id = \"user_page\".user_id WHERE \"isActive\" = true AND \"isPublic\" = true ORDER BY subscribers DESC, user_name")

	if err != nil {
		util.LogError(err, "database")
	}
	defer rows.Close()

	// load data into PageData host list
	for rows.Next() {
		rows.Scan(&name, &url, &avatar, &subs, &events)
		data.HostList = append(data.HostList, models.User{Name: name, Url: url, Avatar: avatar, Subs: subs, Events: events})
	}

	// render
	err = util.RenderTemplate(res, "hosts", false, data)
	if err != nil {
		util.LogError(err, "render")
	}

}
