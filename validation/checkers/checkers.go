package checkers

import (
	// Project Libs
	patterns "hermes-messaging-service/validation/patterns"
)

// IsTokenValid : Checks if parameter matches regex
func IsTokenValid(s string) bool {
	return patterns.RegexToken.MatchString(s)
}
