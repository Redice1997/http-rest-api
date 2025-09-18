package storage

import (
	"context"

	"github.com/Redice1997/http-rest-api/internal/app/model"
)

// UserRepository defines methods for user data management.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}
