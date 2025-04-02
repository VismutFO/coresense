package utils

import (
	"errors"
)

var (
	ErrDuplicatedKey = errors.New("unique key constraint violation")
)
