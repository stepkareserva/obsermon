package dbstorage

import (
	"errors"
	"net"
)

type NetError struct{}

var _ error = (*NetError)(nil)

func (e NetError) Error() string {
	return "temporary net error"
}

func (e NetError) Is(target error) bool {
	var netErr net.Error
	return errors.As(target, &netErr) && netErr.Timeout()
}

var ErrNet NetError
