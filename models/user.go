package models

// Holds user data from the database.
type User struct {
	Name   string
	Url    string
	Avatar string
	Active bool
	Subs   int
}
