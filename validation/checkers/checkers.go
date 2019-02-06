package checkers

import (
	regexp "regexp"
	models "wave-messaging-management-service/models"
)

// IsTokenValid : Checks if parameter matches regex
func IsTokenValid(env *models.Env, s string) (bool, error) {

	// Refresh config to get actual environment values
	err := env.RefreshConfig()

	if err != nil {
		return false, err
	}

	// RegexToken : Regex validation for token (Provided by use through environment variable)
	RegexToken := regexp.MustCompile(env.Config.TokenValidationRegex)

	return RegexToken.MatchString(s), nil
}
