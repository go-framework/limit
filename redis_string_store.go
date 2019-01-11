package limit

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/go-framework/errors"
)

// Redis string impl limit store.
type RedisStringStore struct {
	*redis.Client
}

// Get the identifier key value from Store.
func (this *RedisStringStore) Get(key string) (int, error) {
	count, err := this.Client.Get(key).Int()
	if err == redis.Nil {
		return 0, errors.NewCode(ErrNoStoreKey, "no exist store key "+key)
	} else if err != nil {
		return 0, errors.NewCode(ErrStoreOperate, err.Error())
	}

	return count, nil
}

// Increase or Decrease the identifier key value by count, duration is
// validity period of identifier, when 0 is never expired.
// return latest value.
func (this *RedisStringStore) Set(key string, count int, duration time.Duration) (int, error) {
	// lua script
	// incrby key by count
	// expire key duration'second
	luaScript := `
	local count = redis.call('INCRBY', KEYS[1], ARGV[1])
	local ttl = redis.call('TTL', KEYS[1])
	if ttl <= 0 then
		redis.call('EXPIRE', KEYS[1], ARGV[2])
	end

	return count
`
	return this.Eval(luaScript, []string{key}, count, duration.Seconds()).Int()
}

// New redis string store with redis option.
func NewRedisStringStore(option *redis.Options) Store {
	// default option.
	if option == nil {
		option = &redis.Options{
			Addr: "127.0.0.1:6379",
		}
	}

	object := &RedisStringStore{
		Client: redis.NewClient(option),
	}

	return object
}
