package util

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Get a random string of chars, about half are numbers and the other half are characters
//
// Will give x groups of n chars, separated by a dash.
//
// Return error if n or x is not above 0
func RandomString(x int, n int) (string, error) {

	if x <= 0 || n <= 0 {
		return "", errors.New("RandomString() : " + "x and n must be above 0")
	}

	output := ""

	rand.NewSource(time.Now().UnixMicro())

	for i := 0; i < x; i++ {
		for a := 0; a < n; a++ {

			// do char instead?
			if int(rand.Intn(10)) <= 4 {

				newChar := string(rune(rand.Intn(123-97) + 97))

				// uppercase>
				if int(rand.Intn(10)) >= 5 {
					newChar = strings.ToUpper(newChar)
				}

				output += newChar

			} else {
				output += fmt.Sprintf("%d", rand.Intn(10))
			}

		}

		if i < x-1 {
			output += "-"
		}
	}

	return output, nil

}
