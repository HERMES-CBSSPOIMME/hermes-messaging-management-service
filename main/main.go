package main

import (
	// Native Go Libs
	fmt "fmt"

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

	// RedisHost : Redis Port
	RedisHost = "localhost"

	// RedisPort : Redis Port
	RedisPort = 6379

	// RedisPassword : Redis Password
	RedisPassword = "example"

	// RedisURL : Redis Connection URL
	RedisURL = fmt.Sprintf("redis://%s:%d", RedisHost, RedisPort)
)

func main() {

	// Get MongoDB communication interface
	// If an error occurs, program is set to panic
	mongoDB := models.NewMongoDB(MongoDBURL)

	// Get Redis communication interface
	// If an error occurs, program is set to panic
	redis := models.NewRedis(RedisURL, RedisPassword)

	// Add interfaces & blank config to the environment
	env := &models.Env{
		MongoDB: mongoDB,
		Redis:   redis,
		Config:  models.Config{},
	}

	// Dynamically load config
	env.RefreshConfig()

	router.Listen(env)

	defer func() {
		env.Redis.CloseConnection()
	}()
}
