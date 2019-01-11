package store

import (
	"github.com/go-redis/redis"

	"github.com/go-framework/limit"
)

// New redis string store with redis option.
func NewRedisStringStore(option *redis.Options) limit.Store {
	// default option
	if option == nil {
		option = &redis.Options{
			Addr: "127.0.0.1:6379",
		}
	}

	object := &limit.RedisStringStore{
		Client: redis.NewClient(option),
	}

	return object
}
