package models

import (
	redis "github.com/gomodule/redigo/redis"
)

// RedisInterface : Redis Communication interface
type RedisInterface interface {
	CloseConnection() error
}

// Redis : Redis communication interface
type Redis struct {
	Connection redis.Conn
}

// NewRedis : Return a new Redis abstraction struct
func NewRedis(connectionURL string) *Redis {

	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL(connectionURL)
	if err != nil {
		panic(err)
	}

	conn.Do("AUTH")
	// Return new MongoDB abstraction struct
	return &Redis{
		Connection: conn,
	}
}

// CloseConnection : Close Redis Connection
func (redis *Redis) CloseConnection() error {

	return redis.Connection.Close()
}
