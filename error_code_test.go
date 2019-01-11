package limit

import "testing"

func TestErrCode_Message(t *testing.T) {
	ds := []ErrCode{
		ErrStoreOperate,
		ErrNoStoreKey,
		ErrCheckLimitFailed,
	}

	for i, d := range ds {
		t.Log(i, d)
	}
}
