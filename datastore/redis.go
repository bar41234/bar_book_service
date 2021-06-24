package datastore

import (
	"fmt"
	"gopkg.in/redis.v5"
)

const (
	actionsNumber = 3
)

type redisUserActivity struct {
	client *redis.Client
}

func NewRedisUserActivity(url string) (*redisUserActivity, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &redisUserActivity{client: client}, nil

}

func (e *redisUserActivity) GetActivities(username string) ([]string, error) {
	ret, err := e.client.LRange(username, 0, actionsNumber).Result()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *redisUserActivity) AddActivity(username string, methodName string, request string) {
	if username == "" {
		return
	}
	e.client.LPush(username, fmt.Sprintf("Method: %s, Route: %s", methodName, request))
	length := e.client.LLen(username).Val()
	for length > actionsNumber {
		e.client.RPop(username)
		length--
	}

}
