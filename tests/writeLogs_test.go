package tests

import (
	"errors"
	"outatime/util"
	"testing"
)

// Test that an false is returned if bad file name is given
func TestLogErrorBadLogName(t *testing.T) {

	logged := util.LogError(errors.New("error"), "Bad_file")

	// if true, throw error
	if logged {
		t.Errorf("Bad_file.log does not exist.")
	}

}

// Test that an true is returned if good file name is given
func TestLogErrorGoodLogName(t *testing.T) {

	logged := util.LogError(errors.New("error"), "test")

	// if not true, throw error
	if !logged {
		t.Errorf("test.log does exist.")
	}

}
