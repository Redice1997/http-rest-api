package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/app"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int64
	Email             string
	Password          string
	EncryptedPassword string
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

var (
	ErrValidation    = errors.New("invalid email or password")
	ErrAlreadyExists = errors.New("already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNotFound      = errors.New("not found")
)

func New(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}

func Me(ctx context.Context, r UserRepository, ID int64) (*User, error) {
	return r.GetByID(ctx, ID)
}

func (u *User) Save(ctx context.Context, r UserRepository) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if _, err := r.GetByEmail(ctx, u.Email); errors.Is(err, ErrNotFound) {
	} else if err != nil {
		return err
	} else {
		return ErrAlreadyExists
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if err := r.Create(ctx, u); err != nil {
		return err
	}

	return nil
}

func (u *User) MakeSession(ctx context.Context, ur UserRepository, s sessions.Store, w http.ResponseWriter, r *http.Request) error {
	
	if user, err := ur.GetByEmail(ctx, u.Email); errors.Is(err, ErrNotFound) {
		return ErrValidation
	} else if err != nil {
		return err
	} else if !u.ComparePassword(user.Password) {
		return ErrValidation
	}

	session, err := s.Get(r, app.SessionName)
	if err != nil {
		return err
	}

	session.Values[app.SessonUserID] = u.ID

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (u *User) Validate() error {
	if err := validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	); err != nil {
		return fmt.Errorf("%w: %w", ErrValidation, err)
	}
	return nil
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		encryptedPassword, err := encryptPassword(u.Password)
		if err != nil {
			return fmt.Errorf("before create: %w", err)
		}
		u.EncryptedPassword = encryptedPassword
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func requiredIf(cond bool) validation.RuleFunc {
	return func(value any) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}
