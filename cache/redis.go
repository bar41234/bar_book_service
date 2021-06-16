package cache

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
)

const (
	SERVER_URL = "localhost:6379"
)

type RedisCache struct {
	client *redis.Client
	isInit bool
}

func (e *RedisCache) Init() {
	if !e.isInit {
		client := redis.NewClient(&redis.Options{
			Addr:     SERVER_URL,
			Password: "",
			DB:       0,
		})
		e.client = client
		e.isInit = true
	}
}

func (e *RedisCache) CacheGet(username string) ([]string, error) {
	e.Init()
	ret , err := e.client.LRange(username, 0, 5).Result()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *RedisCache) CacheAdd(username string, ctx *gin.Context) {
	e.Init()
	if username != "" {
		length := e.client.LLen(username)
		if length.Val() >= 6 {
			e.client.RPop(username)
			e.client.RPop(username)
		}
		e.client.LPush(username, ctx.Request.RequestURI)
		e.client.LPush(username, ctx.Request.Method)
	}
}
