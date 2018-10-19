package utils

import (
	log "log"
)

// MappingRequestBody : Request Body on Mapping Request
type MappingRequestBody struct {
	UserIDs []string `json:"userIDs"`
}

// GroupConversationBody : Request Body on Group Creation
type GroupConversationBody struct {
	Members []string `json:"members"`
	Name    string   `json:"name"`
}

// AuthCheckerBody : Response Body from Auth Checker
type AuthCheckerBody struct {
	OriginalUserID string `json:"userID" bson:"userID"`
}

// PanicOnError : Prints the error & exits the program
func PanicOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
