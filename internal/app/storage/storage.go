package storage

// Storeage is the interface that wraps the User method.
type Storage interface {
	User() UserRepository
}
