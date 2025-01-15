package testcommon

import (
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

//var redisOnce sync.Once

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func GetRedisClient() *redis.Client {
	return redisClient
}
