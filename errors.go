package tsl2591

import (
	"errors"
	"fmt"
)

var ErrOverflow = errors.New("overflow reading light channels")

type UnexpectedDeviceIDError struct {
	Expected byte
	Actual   byte
}

func (e UnexpectedDeviceIDError) Error() string {
	return fmt.Sprintf("received device ID %x does not match expected device ID %x", e.Actual, e.Expected)
}
