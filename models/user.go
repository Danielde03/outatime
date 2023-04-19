package models

// Holds user data from the database.
type User struct {
	Name   string
	Url    string
	Avatar string
	Active bool
	Subs   int
	Events int
}

type UserPage struct {
	About  string
	Banner string
	Public bool
}
