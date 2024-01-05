package util

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"outatime/models"
	"strings"
	"text/template"
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

// Render the layout and the chosen template
//
// page is the name of the page template and pass the data to be used in the page in pageData.
//
// logged in to see what nav bar to use. If not, use standard nav bar.
//
// Return error if an error if one appears in parsing or executing templates.
func RenderTemplate(res http.ResponseWriter, page string, loggedIn bool, pageData models.PageData) error {

	t, err := template.ParseFiles("templates/layout.html", "templates/include/"+page+".html")

	// parse nav bar
	if loggedIn {
		_, err := t.ParseFiles("templates/include/navbars/userNav.html")

		if err != nil {
			return err
		}
	} else {
		_, err := t.ParseFiles("templates/include/navbars/stndNav.html")

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	err = t.Execute(res, pageData)

	if err != nil {
		return err
	}

	return nil
}

// list of log files
var loglist = map[string]bool{
	"database": true,
	"files":    true,
	"render":   true,
	"util":     true,
	"overflow": true,
	"test":     true,
	"cookies":  true,
}

// Write errors to a specific log file
//
// err is the error being written, and logName is the log file, excluding the .log extention
//
// If the error was logged in a valid log file, return true.
// If the error was loggen in an invalid log file, return false.
//
// Errors sent to invalid logs will be sent to the overflow.log, and the entry of bad data will be logged in the util.log file.
func LogError(err error, logName string) bool {

	if !loglist[logName] {

		// log bad data in util.log
		utilLog, _ := os.OpenFile("logs/util.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(utilLog)
		log.Println("LogError() : " + logName + " is not a valid log name")

		// store err in overflow
		overflowLog, _ := os.OpenFile("logs/overflow.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(overflowLog)
		log.Println("Invalid log name was given. Error: ", err)

		defer utilLog.Close()
		defer overflowLog.Close()

		return false
	}

	// log error
	fileName := "logs/" + logName + ".log"
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(file)
	log.Println(err)

	defer file.Close()

	return true
}

// Check that the toekn cookie and the token form value match
//
// If not, send 500 error and return false.
//
// If the token and cookie are good, return true.
func CheckToken(res http.ResponseWriter, req *http.Request) bool {

	// get token data
	token := req.FormValue("token")
	tokenCookie, err := req.Cookie("token")

	// remove token cookie
	http.SetCookie(res, &http.Cookie{
		Name:   "token",
		Path:   "/",
		MaxAge: -1,
	})

	// if no token - sign of bad intent
	if err != nil {
		LogError(err, "cookies")
		http.Error(res, "Token cookie not found", 500)
		return false
	}

	// if token don't match - sign of bad intent
	if token != tokenCookie.Value || strings.TrimSpace(token) == "" {
		http.Error(res, "Token does not match", 500)
		return false
	}

	return true

}
