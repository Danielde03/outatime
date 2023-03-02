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
