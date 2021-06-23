package datastore

import (
	"gopkg.in/redis.v5"
)

const (
	actionsNumber = 3
)

type RedisDataStore struct {
	client *redis.Client
}

func NewRedis(url string) (*RedisDataStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisDataStore{client: client}, nil

}

func (e *RedisDataStore) Get(username string) ([]string, error) {
	ret, err := e.client.LRange(username, 0, actionsNumber).Result()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *RedisDataStore) Add(username string, methodName string, request string) {
	if username != "" {
		length := e.client.LLen(username).Val()
		for length >= actionsNumber {
			e.client.RPop(username)
			length--
		}
		e.client.LPush(username, methodName+request)
	}
}
