package sqlstorage

import (
	"context"

	"github.com/Redice1997/http-rest-api/internal/app/model"
)

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// Implement the logic to get a user by email from the SQL database
	return nil, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// Implement the logic to get a user by ID from the SQL database
	return nil, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	// Implement the logic to update a user in the SQL database
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	// Implement the logic to delete a user from the SQL database
	return nil
}
