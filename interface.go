package limit

import (
	"time"
)

// Options map, key type is OptionKey.
type Options map[OptionKey]interface{}

// Option type.
type Option func(Options)

// Limit check and renew pair interface.
type Limit interface {
	// Check if the identifier exceeds the limit, return nil check succeed,
	// otherwise return ErrCheckLimitFailed code error.
	Check(identifier string, opts ...Option) error
	// Renew identifier value with option, return latest value.
	Renew(identifier string, opts ...Option) (int, error)
}

// Limit store interface.
type Store interface {
	// Get the identifier key value from Store.
	// if no such key return ErrNoStoreKey code error.
	Get(key string) (int, error)
	// Increase or Decrease the identifier key value by count, duration is
	// validity period of identifier, when 0 is never expired.
	// return latest value.
	Set(key string, count int, duration time.Duration) (int, error)
}
