package models

import (
	json "encoding/json"
	ioutil "io/ioutil"
	os "os"
)

var (
	configFilePath = os.Getenv("WAVE_CONFIG_FILE_PATH")
)

// Env : Execution environment containing Datastore communication interfaces (Redis, MongoDB) & Config
type Env struct {
	MongoDB MongoDBInterface
	Redis   RedisInterface
	Config  Config
}

// Config : Global Config
type Config struct {
	AuthenticationCheckEndpoint string `json:"authenticationCheckEndpoint"`
	TokenValidationRegex        string `json:"tokenValidationRegex"`
}

// RefreshConfig : Load current environment values in config
func (env *Env) RefreshConfig() error {

	data, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &env.Config)

	if err != nil {
		return err
	}

	return nil
}
