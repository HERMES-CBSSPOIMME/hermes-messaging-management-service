package logruswrapper

import (
	uuid "github.com/satori/go.uuid"
)

const (
	CodeSuccess      = "SUCCESS"
	CodeBadLogin     = "BAD-LOGIN"
	CodeInvalidToken = "INVALID-TOKEN"
)

var (
	messageCodeMapping = map[string]string{
		CodeBadLogin:     "Bad username or password",
		CodeInvalidToken: "Invalid Token",
		CodeSuccess:      "Success",
	}
)

// LogEntryInfos informations about log entry
type LogEntryInfos struct {
	LogID    string `json:"logID"`
	Session  string `json:"session,omitempty"`
	Service  string `json:"service"`
	Endpoint string `json:"endpoint"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

// NewEntryWithSession return a LogEntryInfos struct instance with a session
func NewEntryWithSession(session string, service string, endpoint string, code string) *LogEntryInfos {

	uuid := uuid.NewV4()

	return &LogEntryInfos{
		LogID:    uuid.String(),
		Session:  session,
		Service:  service,
		Endpoint: endpoint,
		Code:     code,
		Message:  messageCodeMapping[code],
	}
}

// NewEntry return a LogEntryInfos struct instance
func NewEntry(service string, endpoint string, code string) *LogEntryInfos {

	uuid := uuid.NewV4()

	return &LogEntryInfos{
		LogID:    uuid.String(),
		Service:  service,
		Endpoint: endpoint,
		Code:     code,
		Message:  messageCodeMapping[code],
	}
}
