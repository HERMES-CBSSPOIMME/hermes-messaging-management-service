package main

import (
	fmt "fmt"
	models "hermes-messaging-management-service/models"
	router "hermes-messaging-management-service/router"
	log "log"
	os "os"
)

var (

	// TODO: Change these to be fetched automatically with Kubernetes Secrets

	// MongoDBHost : MongoDB Host
	MongoDBHost = "mongodb"

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
	RedisHost = "redis-sessions-cache"

	// RedisPort : Redis Port
	RedisPort = 6379

	// RedisPassword : Redis Password
	RedisPassword = "example"

	// RedisURL : Redis Connection URL
	RedisURL = fmt.Sprintf("redis://%s:%d", RedisHost, RedisPort)
)

func main() {

	if os.Getenv("HERMES_CONFIG_FILE_PATH") == "" {
		log.Fatalf("HERMES_CONFIG_FILE_PATH Environment variable must be set !")
	}

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
	err := env.RefreshConfig()

	if err != nil {
		log.Fatalf(err.Error())
	}

	router.Listen(env)

	defer func() {
		env.Redis.CloseConnection()
	}()
}
