package limit

// Error code type.
type ErrCode int32

// Error code defined.
const (
	ErrStoreOperate     ErrCode = -1 * (iota + 1) // store operate error.
	ErrNoStoreKey                                 // no such key in store.
	ErrCheckLimitFailed                           // check limit failed.
)

// Error code messages.
var ErrCodeMessages = []string{
	(-1 * ErrStoreOperate) - 1:     "Store operate error",
	(-1 * ErrNoStoreKey) - 1:       "Store key not exist",
	(-1 * ErrCheckLimitFailed) - 1: "Check limit failed",
}

// Implement Message interface.
func (e ErrCode) Message() string {
	return ErrCodeMessages[(-1*e)-1]
}

// Implement Message interface.
func (e ErrCode) Error() string {
	return ErrCodeMessages[(-1*e)-1]
}
