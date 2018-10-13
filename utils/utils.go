package utils

import (
	log "log"
	"net/http"
)

type (
	Handler func(w http.ResponseWriter, r *http.Request) error
)

// PanicOnError : Prints the error & exits the program
func PanicOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
