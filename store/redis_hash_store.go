package store

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/go-framework/errors"
	"github.com/go-framework/limit"
)

// Redis hash impl limit store.
type RedisHashStore struct {
	*redis.Client
	key string // hash key.
}

// Get the identifier key filed value from Store.
func (this *RedisHashStore) Get(field string) (int, error) {
	count, err := this.Client.HGet(this.key, field).Int()
	if err == redis.Nil {
		return 0, errors.NewCode(limit.ErrNoStoreKey, "no exist store key "+this.key)
	} else if err != nil {
		return 0, errors.NewCode(limit.ErrStoreOperate, err.Error())
	}

	return count, nil
}

// Increase or Decrease the identifier key field value by count, duration is
// validity period of identifier, when 0 is never expired.
// return latest value.
func (this *RedisHashStore) Set(field string, count int, duration time.Duration) (int, error) {
	// lua script
	// hincrby key field by count
	// expire key duration'second
	luaScript := `
	local count = redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2])
	local ttl = redis.call('TTL', KEYS[1])
	if ttl <= 0 then
		redis.call('EXPIRE', KEYS[1], ARGV[3])
	end

	return count
`
	return this.Eval(luaScript, []string{this.key}, field, count, duration.Seconds()).Int()
}

// New redis hash store with store key and redis option.
func NewRedisHashStore(key string, option *redis.Options) limit.Store {
	// default option
	if option == nil {
		option = &redis.Options{
			Addr: "127.0.0.1:6379",
		}
	}

	object := &RedisHashStore{
		Client: redis.NewClient(option),
		key:    key,
	}

	return object
}
