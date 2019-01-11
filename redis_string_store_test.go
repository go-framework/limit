package limit

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"

	"github.com/go-framework/errors"
)

func TestRedisStringStore_Set(t *testing.T) {
	// use mini redis
	mini, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	// address
	addr := mini.Addr()
	// addr = "127.0.0.1:6379"
	// option
	redisOption := &redis.Options{
		Addr: addr,
	}
	// new
	store := NewRedisStringStore(redisOption)

	// params
	key := "TestRedisStringStore"
	count := 1
	duration := time.Second * 30

	// set
	result, err := store.Set(key, count, duration)
	if err != nil {
		t.Fatal(err)
	}
	if result != count {
		t.Errorf("set result should be %d, not %d", count, result)
	}

	if !mini.Exists(key) {
		t.Errorf("%s should be exist", key)
	}

	mini.FastForward(duration)

	if mini.Exists(key) {
		t.Errorf("%s should not be exist", key)
	}

	count = 10
	// set
	result, err = store.Set(key, count, duration)
	if err != nil {
		t.Fatal(err)
	}
	if result != count {
		t.Errorf("set result should be %d, not %d", count, result)
	}

	count2 := -5
	// set
	result, err = store.Set(key, count2, duration)
	if err != nil {
		t.Fatal(err)
	}
	if result != count+count2 {
		t.Errorf("set result should be %d, not %d", count+count2, result)
	}
}

func TestRedisStringStore_Get(t *testing.T) {
	// use mini redis
	mini, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	// address
	addr := mini.Addr()
	// addr = "127.0.0.1:6379"
	// option
	redisOption := &redis.Options{
		Addr: addr,
	}
	// new
	store := NewRedisStringStore(redisOption)

	// params
	key := "TestRedisStringStore"
	count := 1
	duration := time.Second * 30

	// set
	result, err := store.Set(key, count, duration)
	if err != nil {
		t.Fatal(err)
	}
	if result != count {
		t.Errorf("set result should be %d, not %d", count, result)
	}
	// get
	value, err := store.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if value != result {
		t.Errorf("set result should be %d, not %d", result, value)
	}

	mini.FastForward(duration)

	// get
	value, err = store.Get(key)
	if !errors.CodeEqual(ErrNoStoreKey, err) {
		t.Fatal(err)
	}
}
