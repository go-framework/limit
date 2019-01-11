package limit

import (
	"time"

	"github.com/go-framework/errors"
)

// Multiple count limit key prefix.
const MultiCountLimitKeyPrefix = "limit:multi"

// multiple max count limit, impl Limit interface.
type multiMaxCountLimit struct {
	store            Store
	durationMaxCount map[time.Duration]int // key is duration, value is max count.
	keyPrefix        string
}

// Get unique store key.
func (this *multiMaxCountLimit) StoreKey(identifier string, duration time.Duration) string {
	return this.keyPrefix + identifier + duration.String()
}

// Check identifier is exceed the max count limit, return nil is not exceed.
func (this *multiMaxCountLimit) Check(identifier string, opts ...Option) error {

	for duration, max := range this.durationMaxCount {
		value, err := this.store.Get(this.StoreKey(identifier, duration))
		if err != nil {
			return err
		}
		if value > max {
			return errors.NewCodeSprintf(ErrCheckLimitFailed, "check %s value is %d more than %d in %v", identifier, value, max, duration)
		}
	}

	return nil
}

// Renew identifier value by count.
// return latest value in option.
func (this *multiMaxCountLimit) Renew(identifier string, opts ...Option) (int, error) {
	durationMaxCount := make(map[time.Duration]int)
	option := Options{}

	for _, opt := range opts {
		opt(option)
	}

	if value, ok := option[ValueKey]; ok {
		if value2, ok := value.(map[time.Duration]int); ok {
			durationMaxCount = value2
		}
	}

	for duration := range this.durationMaxCount {
		count := 1
		// get Renew count
		if value, ok := durationMaxCount[duration]; ok {
			count = value
		}
		// update value
		result, err := this.store.Set(this.StoreKey(identifier, duration), count, duration)
		if err != nil {
			return 0, err
		}
		// return set result in durationMaxCount
		durationMaxCount[duration] = result
	}

	return 0, nil
}

// option of multiMaxCountLimit to update.
func (this *multiMaxCountLimit) option(opts ...Option) {
	option := Options{}

	for _, opt := range opts {
		opt(option)
	}

	if value, ok := option[StoreKey]; ok {
		if store, ok := value.(Store); ok {
			this.store = store
		}
	}

	if value, ok := option[KeyPrefixKey]; ok {
		if keyPrefix, ok := value.(string); ok {
			this.keyPrefix = keyPrefix
		}
	}

	if value, ok := option[ValueKey]; ok {
		if durationMaxCount, ok := value.(map[time.Duration]int); ok {
			this.durationMaxCount = durationMaxCount
		}
	}
}

// Set multiple count value option.
func WithMultiCountValue(value map[time.Duration]int) Option {
	return func(options Options) {
		options[ValueKey] = value
	}
}

// New multiple max count limit with options.
func NewMultiCountLimit(value map[time.Duration]int, opts ...Option) Limit {
	object := &multiMaxCountLimit{
		store:            DefaultStore,
		durationMaxCount: value,
		keyPrefix:        MultiCountLimitKeyPrefix,
	}
	// refresh.
	object.option(opts...)

	return object
}
