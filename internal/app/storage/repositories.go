package storage

import (
	"context"

	"github.com/Redice1997/http-rest-api/internal/app/models"
)

// UserRepository defines methods for user data management.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}
