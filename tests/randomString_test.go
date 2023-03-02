package tests

import (
	"fmt"
	"math/rand"
	"outatime/util"
	"testing"
	"time"
)

// Is n at or below zero - return error
func TestRandomStringNBelowZero(t *testing.T) {
	_, err := util.RandomString(1, 0)

	if err == nil {

		t.Errorf("Should return error")
	}

}

// Is x at or below 0  - return error
func TestRandomStringXBelowZero(t *testing.T) {
	_, err := util.RandomString(0, 1)

	if err == nil {

		t.Errorf("Should return error")
	}

}

// Test returns proper string length
func TestRandomStringLength(t *testing.T) {

	rand.NewSource(time.Now().UnixMicro())

	x := int(rand.Float64()*10) + 1
	n := int(rand.Float64()*10) + 1

	expected := (n * x) + (x - 1)

	val, _ := util.RandomString(x, n)

	if len(val) != expected {
		t.Errorf("Length should be " + fmt.Sprintf("%d", expected) + " but was " + fmt.Sprintf("%d", len(val)))
	}

}

// Test no repeats. Run 100_000 times to make sure
func TestRandomStringDifferent(t *testing.T) {

	x := 5
	n := 5
	val1, _ := util.RandomString(x, n)

	// run many times
	for i := 0; i < 100_000; i++ {

		val2, _ := util.RandomString(x, n)

		if val1 == val2 {
			t.Errorf("Same string returned twice")
		}

	}

}
