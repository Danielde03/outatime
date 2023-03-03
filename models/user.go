package models

import (
	"fmt"
)

// Holds user data from the database.
type User struct {
	Name   string
	Avatar string
	Active bool
	Subs   int
}

func (u *User) Display() string {
	return fmt.Sprintln(u.Name, u.Active, u.Subs, u.Avatar)
}
