package storage

import "context"

// Storage is the interface that wraps the User method.
type BaseStorage interface {
	User() UserRepository
}

type Storage interface {
	BaseStorage
	BeginTx(ctx context.Context) (TxStorage, error)
}

type TxStorage interface {
	BaseStorage
	Commit() error
	Rollback() error
}
