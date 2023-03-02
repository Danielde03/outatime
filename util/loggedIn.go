package util

import "net/http"

// Determine if a user is logged in.
func IsLoggedIn(req *http.Request) (bool, string) {

	id := req.Header.Get("user_id")

	if len(id) > 0 {
		return true, id
	}

	return false, id

}
