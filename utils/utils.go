package utils

import (
	// Native Go Libs
	log "log"
)

type GroupConversationBody struct {
	Members []string `json="members"`
	Name    string   `json:"name"`
}

// PanicOnError : Prints the error & exits the program
func PanicOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
