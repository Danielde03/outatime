package util

import (
	"log"
	"os"
)

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

		return false
	}

	// log error
	fileName := "logs/" + logName + ".log"
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(file)
	log.Println(err)

	return true
}
