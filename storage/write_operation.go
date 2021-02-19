package storage

import (
	"errors"
)

const (

	// Insert indicates an insertion write operation. It will fail if the target already exists.
	Insert WriteOperation = iota

	// Update indicates an update write operation. It will fail if the target does not exist.
	Update

	// Upsert indicates an upsert write operation. It will be performed if the target exists or does not exist.
	Upsert
)

var (

	// ErrNotWriteOperation indicates that the given input string was not a known write operation.
	ErrNotWriteOperation = errors.New("the input string was not a known write operation")
)

// WriteOperation indicates the type of write operation.
type WriteOperation uint

// FromString will determine what write operation should take place based on the input string.
func (w WriteOperation) FromString(s string) (operation WriteOperation, err error) {
	switch s {
	case "insert":
		operation = Insert
	case "update":
		operation = Update
	case "upsert":
		operation = Upsert
	default:
		return Insert, ErrNotWriteOperation
	}
	return operation, nil
}
