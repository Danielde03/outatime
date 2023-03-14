package tests

import (
	"fmt"
	"outatime/util"
	"testing"
)

// Test we can access the database
func TestCanAccessDB(t *testing.T) {

	_, err := util.DatabaseExecute("SELECT sub_id FROM outatime.subscriber")

	if err != nil {
		t.Errorf("An error was returned when executing")
	}
}

// Test errors are thrown
func TestErrorThrows(t *testing.T) {

	_, err := util.DatabaseExecute("SELECT sub_id, FROM outatime.subscriber")

	if err == nil {
		t.Errorf("No error was returned when executing")
	}
}

// Make sure unique values return false for isAdmin
//
// DATABASE MUST HAVE d@p.ca TEST ACCOUNT. WILL FAIL IF NOT
func TestCaseSensitive(t *testing.T) {

	_, err := util.DatabaseExecute(fmt.Sprintf("Select \"%v\" From outatime.user where user_email = 'd@p.ca'", "isAdmin"))

	if err != nil {
		t.Errorf(err.Error())
	}
}

// Make sure unique values return true
func TestIsUnique(t *testing.T) {

	out := util.IsUnique("hdffdnmdngdndgnn", "user_email", "user")

	if !out {
		t.Errorf("IsUnique should be true, but returns false")
	}
}

// Make sure unique values return true
//
// DATABASE MUST HAVE d@p.ca TEST ACCOUNT. WILL FAIL IF NOT
func TestIsNotUnique(t *testing.T) {

	out := util.IsUnique("d@p.ca", "user_email", "user")

	if out {
		t.Errorf("IsUnique should be false, but returns true. - Make sure d@p.ca is in database for this test")
	}
}
