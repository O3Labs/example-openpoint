package redis

import (
	"sync"

	"encoding/json"

	"gopkg.in/redis.v5"
)

type Redis struct {
	client *redis.Client
}

var redisClient *Redis
var once sync.Once

func NewRedisClient() *Redis {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		redisClient = &Redis{client: client}
		pong, err := redisClient.client.Ping().Result()
		if err != nil {

			return
		}
		logChannel.Printf("Pong %v", pong)
	})

	return redisClient
}

func (r *Redis) Get(key string) interface{} {
	b, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	return b
}

func (r *Redis) GetScan(key string, v interface{}) error {
	b, err := r.client.Get(key).Bytes()
	if err != nil || b == nil {
		return err
	}
	json.Unmarshal(b, &v)
	return err
}

func (r *Redis) Set(key string, value interface{}) {
	if value == nil {
		r.client.Del(key)
		return
	}
	b, _ := json.Marshal(value)
	r.client.Set(key, b, 0)
}
