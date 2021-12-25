package client_singleton

import (
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
)
