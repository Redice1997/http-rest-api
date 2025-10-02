package sqlstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Redice1997/http-rest-api/internal/app/model"
)

type userRepository struct {
	db sqlDB
}

func newUserRepository(db sqlDB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	return r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {

	u := new(model.User)

	if err := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {

	u := new(model.User)

	if err := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) Update(ctx context.Context, u *model.User) error {

	_, err := r.db.ExecContext(
		ctx,
		"UPDATE users SET email = $1, encrypted_password = $2 WHERE id = $3",
		u.Email,
		u.EncryptedPassword,
		u.ID,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {

	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)

	return err
}
