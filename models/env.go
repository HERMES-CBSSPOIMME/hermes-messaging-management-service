package models

// Env : Execution environment containing RabbitMQ and MongoDB communication interfaces
type Env struct {
	DB DatastoreInterface
}
