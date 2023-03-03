package routes

import (
	"net/http"
	"outatime/handlers"
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

	// page routes
	mux.HandleFunc("/", rootOr404)

	return mux
}

// If no static or dynamic page was found, check if root. If is, send to root handler, else open 404 page
func rootOr404(res http.ResponseWriter, req *http.Request) {

	// Handle dynamic routes. If no page was found, check if 404 or root.
	if !checkForDynamicRoutes(res, req) {

		// Get user logged in status.
		loggedIn, _ := util.IsLoggedIn(req)

		// open 404 page if not root, and not dynamic route
		if req.URL.Path != "/" {
			err := util.RenderTemplate(res, "404", loggedIn, nil)
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

		// Route
		if count > 0 {
			found = routeDynamicURL(url, user_url, res, req)

		}

		defer rows.Close()

	}
	return found

}

// Get rit of empty strings in array
func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// If a route has a dynamic user, route it based on trailing address
//
// Return true if a page was found.
func routeDynamicURL(url string, user_url string, res http.ResponseWriter, req *http.Request) bool {

	// get url after user's url
	restOfUrl := strings.SplitAfter(url, "/"+user_url)[1]
	if restOfUrl == "" {
		restOfUrl = "/"
	}

	// route
	switch restOfUrl {
	case "/":
		handlers.UserHome(res, req)
		return true
	default:
		return false
	}

}
