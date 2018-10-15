package models

import "os"

// Env : Execution environment containing Datastore communication interfaces (Redis, MongoDB) & Config
type Env struct {
	MongoDB MongoDBInterface
	Redis   RedisInterface
	Config  Config
}

// Config : Global Config
type Config struct {
	AuthenticationCheckEndpoint string
	TokenValidationRegex        string
}

// RefreshConfig : Load current environment values in config
func (env *Env) RefreshConfig() {

	env.Config = Config{
		AuthenticationCheckEndpoint: os.Getenv("HERMES_AUTH_CHECK_ENDPOINT"),
		TokenValidationRegex:        os.Getenv("HERMES_TOKEN_VALIDATION_REGEX"),
	}
}
