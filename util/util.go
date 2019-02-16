package util

import (
	log "github.com/sirupsen/logrus"
)

// RaiseFatalErrorIf stops the program and prints the error message
func RaiseFatalErrorIf(err error, msg string)  {
	if err != nil {
		log.WithError(err).Fatal(msg)
	}
}