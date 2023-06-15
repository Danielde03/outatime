package models

// Holds user data from the database.
type User struct {
	Name       string
	Url        string
	Avatar     string
	Active     bool
	Subs       int
	Events     int
	Event_List []Event
}

// Hold data for a user's page
type UserPage struct {
	About  string
	Banner string
	Public bool
}

// Hold event data from the database
type Event struct {
	Name        string
	Tldr        string
	Description string
	Start       string
	End         string
	Location    string
	Image       string
	IsPrivate   bool
	Code        string
}
