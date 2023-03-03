package models

import "fmt"

// Holds user data from the database.
type User struct {
	Name   string
	Avatar string
	Active bool
	Subs   int
}

func (u *User) Display() {
	fmt.Println(u.Name, u.Active, u.Subs, u.Avatar)
}
