package util

// Get a user's URL based on the ID
//
// Empty URL means no user at that ID
func GetUserURL(id string) string {

	rows, err := DatabaseExecute("SELECT user_url FROM outatime.user WHERE user_id = " + id + ";")

	if err != nil {
		LogError(err, "database")
	}

	userURL := ""

	for rows.Next() {
		err := rows.Scan(&userURL)

		if err != nil {
			LogError(err, "database")
		}
	}

	return userURL

}
