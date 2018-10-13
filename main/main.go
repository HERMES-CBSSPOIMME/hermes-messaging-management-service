package main

import (
	// Native Go Libs
	fmt "fmt"
	os "os"

	// Project Libs
	models "hermes-messaging-service/models"
	router "hermes-messaging-service/router"
)

var (

	// TODO: Change these to be fetched automatically with Kubernetes Secrets

	// MongoDBHost : MongoDB Host
	MongoDBHost = "localhost"

	// MongoDBPort : MongoDB Port
	MongoDBPort = 27017

	// MongoDBUsername : MongoDB Username
	MongoDBUsername = "hermes-user"

	// MongoDBPassword : MongoDB Password
	MongoDBPassword = "example"

	// MongoDBName : MongoDB Database Name
	MongoDBName = "hermesDB"

	// MongoDBURL : MongoDB Connection URL
	MongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", MongoDBUsername, MongoDBPassword, MongoDBHost, MongoDBPort, MongoDBName)
)

func main() {

	// Get MongoDB Communication Interface
	// If an error occurs, program is set to panic
	mongoDB := models.NewMongoDB(MongoDBURL)

	// Add MongoDB Interface to the environment
	env := &models.Env{
		DB: mongoDB,
		Config: models.Config{
			AuthenticationCheckEndpoint: os.Getenv("HERMES_AUTH_CHECK_ENDPOINT"),
		},
	}

	router.Listen(env)

	defer func() {
		// Close stuffs here
	}()
}
