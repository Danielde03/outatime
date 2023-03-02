package tests

import (
	"net/http"
	"outatime/util"
	"strings"
	"testing"
)

// Test a user is logged in
func TestIsLoggedIn(t *testing.T) {

	// make request
	req, _ := http.NewRequest("post", "/", strings.NewReader("Hello, Reader!"))

	// add user id value
	req.Header.Add("user_id", "1")

	loggedIn, _ := util.IsLoggedIn(req)

	if !loggedIn {
		t.Errorf("Said not logged in but is.")
	}

}

// Test right id is returned
func TestRightId(t *testing.T) {

	// make request
	req, _ := http.NewRequest("post", "/", strings.NewReader("Hello, Reader!"))

	// add user id value
	req.Header.Add("user_id", "1")

	_, id := util.IsLoggedIn(req)

	if id != "1" {
		t.Errorf("Returned wrong id.")
	}

}

// Test a user is not logged in
func TestIsNotLoggedIn(t *testing.T) {

	// make request
	req, _ := http.NewRequest("post", "/", strings.NewReader("Hello, Reader!"))

	loggedIn, _ := util.IsLoggedIn(req)

	if loggedIn {
		t.Errorf("Said logged in but is not.")
	}
}
