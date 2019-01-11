package limit

import "time"

// Key of option.
type OptionKey string

const (
	StoreKey     OptionKey = "Store"
	DurationKey  OptionKey = "Duration"
	KeyPrefixKey OptionKey = "KeyPrefix"
	CountKey     OptionKey = "Count"
	ValueKey     OptionKey = "Value"
)

// Set store option.
func WithStore(store Store) Option {
	return func(options Options) {
		options[StoreKey] = store
	}
}

// Set duration option.
func WithDuration(duration time.Duration) Option {
	return func(options Options) {
		options[DurationKey] = duration
	}
}

// Set key prefix option.
func WithKeyPrefix(keyPrefix string) Option {
	return func(options Options) {
		options[KeyPrefixKey] = keyPrefix
	}
}

// Set count option.
func WithCount(count int) Option {
	return func(options Options) {
		options[CountKey] = count
	}
}

// Set value option.
func WithValue(value interface{}) Option {
	return func(options Options) {
		options[ValueKey] = value
	}
}
