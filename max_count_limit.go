package limit

import (
	"strconv"
	"time"

	"github.com/go-framework/errors"
)

// Max count limit key prefix.
const MaxCountLimitKeyPrefix = "limit:max"

// Max count limit, impl Limit interface.
type maxCountLimit struct {
	store     Store
	keyPrefix string
	max       int
	duration  time.Duration
}

// Get unique store key.
func (this *maxCountLimit) StoreKey(identifier string) string {
	return this.keyPrefix + identifier + strconv.Itoa(this.max)
}

// Check identifier is exceed the max count limit, return nil is not exceed.
func (this *maxCountLimit) Check(identifier string, opts ...Option) error {
	value, err := this.store.Get(this.StoreKey(identifier))
	if err != nil {
		return err
	}
	if value > this.max {
		return errors.NewCodeSprintf(ErrCheckLimitFailed, "check %s value is %d more than %d", identifier, value, this.max)
	}

	return nil
}

// Renew identifier value by count.
func (this *maxCountLimit) Renew(identifier string, opts ...Option) (int, error) {
	count := 1
	option := Options{}

	for _, opt := range opts {
		opt(option)
	}

	if value, ok := option[CountKey]; ok {
		if v2, ok := value.(int); ok {
			count = v2
		}
	}

	return this.store.Set(this.StoreKey(identifier), count, this.duration)
}

// option of maxcountLimit to update.
func (this *maxCountLimit) option(opts ...Option) {
	option := Options{}

	for _, opt := range opts {
		opt(option)
	}

	if value, ok := option[StoreKey]; ok {
		if store, ok := value.(Store); ok {
			this.store = store
		}
	}

	if value, ok := option[DurationKey]; ok {
		if duration, ok := value.(time.Duration); ok {
			this.duration = duration
		}
	}

	if value, ok := option[KeyPrefixKey]; ok {
		if keyPrefix, ok := value.(string); ok {
			this.keyPrefix = keyPrefix
		}
	}

	if value, ok := option[ValueKey]; ok {
		if max, ok := value.(int); ok {
			this.max = max
		}
	}
}

// Set max value option.
func WithMaxValue(max int) Option {
	return func(options Options) {
		options[ValueKey] = max
	}
}

// New max count limit with options.
func NewCountLimit(max int, opts ...Option) Limit {
	object := &maxCountLimit{
		store:     DefaultStore,
		max:       max,
		keyPrefix: MaxCountLimitKeyPrefix,
	}
	// refresh.
	object.option(opts...)

	return object
}
