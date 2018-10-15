package checkers

import (
	// Project Libs
	models "hermes-messaging-service/models"
	"regexp"
)

// IsTokenValid : Checks if parameter matches regex
func IsTokenValid(env *models.Env, s string) bool {

	// Refresh config to get actual environment values
	env.RefreshConfig()

	// RegexToken : Regex validation for token (Provided by use through environment variable)
	RegexToken := regexp.MustCompile(env.Config.TokenValidationRegex)

	return RegexToken.MatchString(s)
}
