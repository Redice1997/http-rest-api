package storage

import "errors"

var (
	// ErrRecordNotFound indicates that a record was not found in the storage.
	ErrRecordNotFound = errors.New("record not found")
)
