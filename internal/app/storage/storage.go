package storage

import "context"

// Storage is the interface that wraps the User method.
type Storage interface {
	BeginTx(ctx context.Context) (TxStorage, error)
	User() UserRepository
}

type TxStorage interface {
	Commit() error
	Rollback() error
	User() UserRepository
}
