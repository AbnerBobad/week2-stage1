package data

import "errors"

var (
	// ErrRecordNotFound is returned when a database query finds no
	// matching record.
	ErrRecordNotFound = errors.New("Database record not found")
)
