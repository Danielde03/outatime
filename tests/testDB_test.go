package tests

import (
	"outatime/util"
	"testing"
)

// Test we can access the database
func TestCanAccessDB(t *testing.T) {

	_, err := util.DatabaseExecute("SELECT sub_id FROM outatime.subscriber")

	if err != nil {
		t.Errorf("An error was returned when exitng")
	}
}

// Test errors are thrown
func TestErrorThrows(t *testing.T) {

	_, err := util.DatabaseExecute("SELECT sub_id, FROM outatime.subscriber")

	if err == nil {
		t.Errorf("No error was returned when exitng")
	}
}

// Make sure unique values return true
func TestIsUnique(t *testing.T) {

	if !util.IsUnique("hdffdnmdngdndgnn", "user_email", "user") {
		t.Errorf("IsUnique should be true, but returns false")
	}
}

// Make sure unique values return true
//
// DATABASE MUST HAVE d@p.ca TEST ACCOUNT. WILL FAIL IF NOT
func TestIsNotUnique(t *testing.T) {

	if util.IsUnique("d@p.ca", "user_email", "user") {
		t.Errorf("IsUnique should be false, but returns true. - Make sure d@p.ca is in database for this test")
	}
}
