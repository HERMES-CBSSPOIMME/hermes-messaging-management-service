package models

import (
	"fmt"

	redisgo "github.com/gomodule/redigo/redis"
)

// RedisInterface : Redis Communication interface
type RedisInterface interface {
	CloseConnection() error
	Get(key string) ([]byte, error)
	HGet(key string, field string) ([]byte, error)
	HSet(key string, field string, value []byte) error
	Set(key string, value []byte) error
	Exists(key string) (bool, error)
	Delete(key string) error
	GetKeys(pattern string) ([]string, error)
	Incr(counterKey string) (int, error)
	Rename(oldKey string, newKey string) error
}

// Redis : Redis communication interface
type Redis struct {
	Connection redisgo.Conn
}

// NewRedis : Return a new Redis abstraction struct
func NewRedis(connectionURL string, password string) *Redis {

	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redisgo.DialURL(connectionURL)
	if err != nil {
		panic(err)
	}

	// Authenticate to Redis
	conn.Do("AUTH", password)

	// Return new MongoDB abstraction struct
	return &Redis{
		Connection: conn,
	}
}

// CloseConnection : Close Redis Connection
func (redis *Redis) CloseConnection() error {

	return redis.Connection.Close()
}

func (redis *Redis) Get(key string) ([]byte, error) {

	var data []byte
	data, err := redisgo.Bytes(redis.Connection.Do("GET", key))

	if err != nil {
		return nil, fmt.Errorf("error getting key %s : %v", key, err)
	}
	return data, nil
}

func (redis *Redis) HGet(key string, field string) ([]byte, error) {

	var data []byte
	data, err := redisgo.Bytes(redis.Connection.Do("HGET", key, field))

	if err != nil {
		return nil, fmt.Errorf("error getting key %s : %v", key, err)
	}
	return data, nil
}

func (redis *Redis) HSet(key string, field string, value []byte) error {

	_, err := redis.Connection.Do("HSET", key, field, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s : %v", key, v, err)
	}
	return nil
}

func (redis *Redis) Set(key string, value []byte) error {

	_, err := redis.Connection.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s : %v", key, v, err)
	}
	return nil
}

func (redis *Redis) Rename(oldKey string, newKey string) error {

	_, err := redis.Connection.Do("RENAME", oldKey, newKey)
	if err != nil {
		return fmt.Errorf("error renaming key %s to %s : %v", oldKey, newKey, err)
	}
	return nil
}

func (redis *Redis) Exists(key string) (bool, error) {

	ok, err := redisgo.Bool(redis.Connection.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists : %v", key, err)
	}
	return ok, nil
}

func (redis *Redis) Delete(key string) error {

	_, err := redis.Connection.Do("DEL", key)

	if err != nil {
		return err
	}

	return nil
}

func (redis *Redis) GetKeys(pattern string) ([]string, error) {

	iter := 0
	keys := []string{}
	for {
		arr, err := redisgo.Values(redis.Connection.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redisgo.Int(arr[0], nil)
		k, _ := redisgo.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (redis *Redis) Incr(counterKey string) (int, error) {

	return redisgo.Int(redis.Connection.Do("INCR", counterKey))
}
