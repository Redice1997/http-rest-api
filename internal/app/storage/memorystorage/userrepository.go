package memorystorage

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

type UserRepository struct {
	s     *Storage
	users map[int64]*model.User
	m     *sync.RWMutex
	id    atomic.Int64
}

func NewUserRepository(s *Storage) *UserRepository {
	return &UserRepository{
		s:     s,
		m:     new(sync.RWMutex),
		users: make(map[int64]*model.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	r.m.Lock()
	defer r.m.Unlock()

	u.ID = r.id.Add(1)

	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, storage.ErrRecordNotFound
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, storage.ErrRecordNotFound
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.users[u.ID]; !ok {
		return storage.ErrRecordNotFound
	}

	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.users[id]; !ok {
		return storage.ErrRecordNotFound
	}

	delete(r.users, id)
	return nil
}
