package models

// All page data structs will be of type PageData, so the data can be rendered in a template.
type PageData struct {
	// The user who owns the page
	PageUser User

	// The user logged in
	NavUser User

	// User page data
	UserPage UserPage

	// list of users
	HostList []User
}
