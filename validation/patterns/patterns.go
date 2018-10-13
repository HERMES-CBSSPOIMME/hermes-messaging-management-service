package patterns

import (
	// Native Go Libs
	os "os"
	regexp "regexp"
)

var (
	// RegexToken : Regex validation for token (Provided by use through environment variable)
	RegexToken = regexp.MustCompile(os.Getenv("HERMES_TOKEN_VALIDATION_REGEX"))
)
