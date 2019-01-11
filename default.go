package limit

// Global default store is redis store.
var DefaultStore Store = NewRedisStringStore(nil)
