package models

import (
	uuid "github.com/satori/go.uuid"
)

// GroupConversation : Group conversation struct
type GroupConversation struct {
	GroupConversationID string   `json:"GroupConversationID" bson:"groupConversationID"`
	Name                string   `json:"name" bson:"name"`
	Members             []string `json:"members" bson:"members"`
	// TODO: Add message backup support
}

// NewGroupConversation : Return new VerneMQACL struct pointer
func NewGroupConversation(name string, members []string) *GroupConversation {
	return &GroupConversation{
		GroupConversationID: uuid.NewV4().String(),
		Name:                name,
		Members:             members,
	}
}
