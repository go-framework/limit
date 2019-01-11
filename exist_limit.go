package limit

import (
	"time"

	"github.com/go-framework/errors"
)

// Exist limit key prefix.
const ExistLimitKeyPrefix = "limit:exist"

// exist limit, impl Limit interface.
type existLimit struct {
	store     Store
	name      string
	duration  time.Duration
	keyPrefix string
}

// Get unique store key.
func (this *existLimit) StoreKey() string {
	return this.keyPrefix + this.name
}

// Check identifier if exist return nil.
// if CountKey option incoming, if identifier count more than CountKey value,
// then exist return nil.
func (this *existLimit) Check(identifier string, opts ...Option) error {

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

	value, err := this.store.Get(identifier)
	if errors.CodeEqual(ErrNoStoreKey, err) {
		return errors.NewCodeSprintf(ErrCheckLimitFailed, "check %s is not exist %s", this.name, identifier)
	} else if err != nil {
		return err
	}

	if value < count {
		return errors.NewCodeSprintf(ErrCheckLimitFailed, "check %s %s value %d less than %d", this.name, identifier, value, count)
	}

	return nil
}

// Renew identifier value by count.
func (this *existLimit) Renew(identifier string, opts ...Option) (int, error) {
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

	return this.store.Set(identifier, count, this.duration)
}

// option of existLimit to update.
func (this *existLimit) option(opts ...Option) {
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
}

// New exist limit with options.
// if identifier exist then Check return nil.
func NewExistLimit(name string, store Store, opts ...Option) Limit {
	object := &existLimit{
		store:     store,
		name:      name,
		keyPrefix: ExistLimitKeyPrefix,
	}
	// refresh.
	object.option(opts...)

	return object
}
