package routes

import (
	"net/http"
	"outatime/handlers"
	"outatime/models"
	"outatime/util"
	"strings"
)

// Get routes Mux
//
// Handle static routes
func GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// connect to public folder
	fs1 := http.FileServer(http.Dir("./templates/public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs1))

	// form routes
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/signup", handlers.SignUp)

	// page routes

	// check if dynamic, 404 or root
	mux.HandleFunc("/", rootOr404)

	return mux
}

// If a route has a dynamic user, route it based on trailing address
//
// Return true if a page was found.
func routeDynamicURL(url string, user_url string, res http.ResponseWriter, req *http.Request) bool {

	// if not logged in to page's owner and page is private, redirect to root.
	loggedIn, loggedIn_id := util.IsLoggedIn(req)
	page := *util.GetUserPage(util.GetUserId(user_url))

	// if page is private and not logged in
	if !page.Public && !loggedIn {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return false

		// if page is private and logged in user is not owner
	} else if !page.Public && loggedIn_id != util.GetUserId(user_url) {
		util.RedirectToUser(res, req)
	}

	// get trailing url after user's url
	restOfUrl := strings.SplitAfter(url, "/"+user_url)[1]
	if restOfUrl == "" {
		restOfUrl = "/"
	}

	// route
	switch restOfUrl {
	case "/":
		handlers.UserHome(res, req, user_url)
		return true
	default:
		return false
	}

}

// If no static or dynamic page was found, check if root. If is, send to root handler, else open 404 page
func rootOr404(res http.ResponseWriter, req *http.Request) {

	// Handle dynamic routes. If no page was found, check if 404 or root.
	if !checkForDynamicRoutes(res, req) {

		// Get user logged in status.
		loggedIn, id := util.IsLoggedIn(req)
		var data models.PageData

		if loggedIn {
			data.NavUser = *util.GetUserById(id)
		}

		// open 404 page if not root, and not dynamic route
		if req.URL.Path != "/" {
			err := util.RenderTemplate(res, "404", loggedIn, data)
			if err != nil {
				util.LogError(err, "render")
			}
			return
		}

		handlers.Root(res, req)

	}

}

// Handle dynamic routes
//
// If url is in database, send to user page handlers, else continue through root page handler.
//
// Return true if a page was found
func checkForDynamicRoutes(res http.ResponseWriter, req *http.Request) bool {

	found := false
	url := req.URL.String()

	split_url := strings.Split(url, "/")

	split_url = delete_empty(split_url)

	// is it a dynamic URL?
	if len(split_url) > 0 {

		user_url := split_url[0]

		// see if url is in database
		rows, err := util.DatabaseExecute("SELECT user_id FROM outatime.user WHERE user_url = '" + user_url + "';")

		if err != nil {
			util.LogError(err, "database")
		}

		count := 0
		for rows.Next() {
			count += 1
		}

		// If is in database, route
		if count > 0 {
			found = routeDynamicURL(url, user_url, res, req)

		}

		defer rows.Close()

	}
	return found

}

// Get rit of empty strings in slice
func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
