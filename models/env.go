package models

// Env : Execution environment containing RabbitMQ and MongoDB communication interfaces
type Env struct {
	DB     DatastoreInterface
	Config Config
}

// Config : Global Config
type Config struct {
	AuthenticationCheckEndpoint string
}
