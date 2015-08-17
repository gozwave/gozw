package util

import "errors"

func CheckSliceLength(slice []byte, index int) error {
	if slice == nil {
		return errors.New("nil slice")
	}

	if len(slice) <= index {
		return errors.New("slice index out of bounds")
	}

	return nil
}
