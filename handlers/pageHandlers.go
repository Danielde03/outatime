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

// Handle /{{user_url}}/account
func Account(res http.ResponseWriter, req *http.Request, user_url string) {

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

	util.RenderTemplate(res, "account", loggedIn, data)

}

// Handle /{{user_url}}/events
func Events(res http.ResponseWriter, req *http.Request, user_url string) {

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

		data.NavUser = *util.GetUserById(loggedIn_id)

		// Get user's events
		rows, err := util.DatabaseExecute("SELECT event_name, event_tldr, event_descr, event_start, event_end, event_location, event_img, \"isPublic\", event_code FROM outatime.event WHERE user_id=$1;", loggedIn_id)

		if err != nil {
			util.LogError(err, "database")
		}

		var name string
		var tldr string
		var descr string
		var start string
		var end string
		var loc string
		var img string
		var priv bool
		var code string

		// load data into PageData host list
		for rows.Next() {
			rows.Scan(&name, &tldr, &descr, &start, &end, &loc, &img, &priv, &code)
			data.PageUser.Event_List = append(data.PageUser.Event_List, models.Event{Name: name, Tldr: tldr, Description: descr, Start: start, End: end, Location: loc, Image: img, IsPrivate: priv, Code: code})
		}

		// if logged in to own, show editable options
		util.RenderTemplate(res, "events-edit", loggedIn, data)

	} else {

		// Get public events
		rows, err := util.DatabaseExecute("SELECT event_name, event_tldr, event_descr, event_start, event_end, event_location, event_img FROM outatime.event WHERE user_id = $1 AND \"isPublic\" = true;", util.GetUserId(user_url))

		if err != nil {
			util.LogError(err, "database")
		}

		var name string
		var tldr string
		var descr string
		var start string
		var end string
		var loc string
		var img string

		// load data into PageData host list
		for rows.Next() {
			rows.Scan(&name, &tldr, &descr, &start, &end, &loc, &img)
			data.PageUser.Event_List = append(data.PageUser.Event_List, models.Event{Name: name, Tldr: tldr, Description: descr, Start: start, End: end, Location: loc, Image: img})
		}

		// If viewers are accessing page, display list of user's public events
		util.RenderTemplate(res, "events-public", loggedIn, data)

	}

}
