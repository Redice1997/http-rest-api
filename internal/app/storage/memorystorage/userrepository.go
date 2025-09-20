package memorystorage

import (
	"context"

	"github.com/Redice1997/http-rest-api/internal/app/model"
)

type UserRepository struct {
	s     *Storage
	users map[string]*model.User
}

func NewUserRepository(s *Storage) *UserRepository {
	return &UserRepository{
		s:     s,
		users: make(map[string]*model.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	panic("implement me")
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("implement me")
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	panic("implement me")
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	panic("implement me")
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}
