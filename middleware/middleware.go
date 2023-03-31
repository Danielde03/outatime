package middleware

import (
	"net/http"
	"outatime/util"
)

// See if a user is logged in
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// get auth_code cookie
		auth_code_cookie, err := req.Cookie("auth_code")

		// if no token, not logged in
		if err != nil {
			next.ServeHTTP(res, req)
			return
		}

		// if logged in, get ID and put in req Header
		auth_code := auth_code_cookie.Value

		// stop users from writing value of null in cookie
		if auth_code == "null" {
			next.ServeHTTP(res, req)
			return
		}

		// get user id from database based on auth_code
		rows, err := util.DatabaseExecute("SELECT user_id FROM outatime.user WHERE auth_code = $1;", auth_code)

		if err != nil {
			util.LogError(err, "database")
			next.ServeHTTP(res, req)
			return
		}

		var id string

		for rows.Next() {
			err := rows.Scan(&id)

			if err != nil {
				util.LogError(err, "database")
				next.ServeHTTP(res, req)
				return
			}
		}

		defer rows.Close()

		// set id in req Header
		req.Header.Add("user_id", id)

		// continue
		next.ServeHTTP(res, req)
	})
}
