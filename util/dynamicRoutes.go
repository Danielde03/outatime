package util

import (
	"fmt"
	"net/http"
	"strings"
)

// Handle dynamic routes
//
// If url is in database, send to user page handlers, else continue through root page handler.
func DynamicRoutes(res http.ResponseWriter, req *http.Request) {

	url := req.URL.String()

	split_url := strings.Split(url, "/")

	split_url = delete_empty(split_url)

	// is it a dynamic URL?
	if len(split_url) > 0 {

		user_url := split_url[0]

		// see if url is in database
		rows, err := DatabaseExecute("SELECT user_id FROM outatime.user WHERE user_url = '" + user_url + "';")

		if err != nil {
			LogError(err, "database")
		}

		count := 0
		for rows.Next() {
			count += 1
		}

		// TODO: if returned value, redirect to user page handler
		if count > 0 {
			fmt.Println("In database")
		}

		defer rows.Close()

	}

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
