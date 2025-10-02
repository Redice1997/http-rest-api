package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	ss sessions.Store
	db storage.Storage
	lg *slog.Logger
}

func New(db storage.Storage, ss sessions.Store, lg *slog.Logger) *Service {
	return &Service{
		db: db,
		ss: ss,
		lg: lg,
	}
}

func Me(ctx context.Context) (*model.User, error) {
	u, ok := ctx.Value(model.CtxUserKey).(*model.User)
	if !ok || u == nil {
		return nil, ErrUserUnauthorized
	}
	return u, nil
}

func (s *Service) Authenticate(r *http.Request) (context.Context, error) {
	session, err := s.ss.Get(r, model.SessionName)
	if err != nil {
		return nil, ErrUserUnauthorized
	}
	id, ok := session.Values[model.SessionUserID].(int64)
	if !ok {
		return nil, ErrUserUnauthorized
	}

	ctx := r.Context()

	u, err := s.db.User().GetByID(ctx, id)
	if errors.Is(err, model.ErrRecordNotFound) {
		return nil, ErrUserUnauthorized
	} else if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, model.CtxUserKey, u), nil
}

func (s *Service) Save(ctx context.Context, u *model.User) error {

	if err := validate(u); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func(tx storage.TxStorage) {
		err := tx.Rollback()
		if err != nil {
			s.lg.Error("failed to rollback transaction", "error", err)
		}
	}(tx)

	if _, err := tx.User().GetByEmail(ctx, u.Email); errors.Is(err, model.ErrRecordNotFound) {

	} else if err != nil {
		return err
	} else {
		return ErrUserExists
	}

	if err := beforeSave(u); err != nil {
		return err
	}

	if err := tx.User().Create(ctx, u); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) CreateHttpSession(w http.ResponseWriter, r *http.Request, u *model.User) error {

	if found, err := s.db.User().GetByEmail(r.Context(), u.Email); errors.Is(err, model.ErrRecordNotFound) {
		return ErrInvalidEmailOrPassword
	} else if err != nil {
		return err
	} else if !comparePassword(found.EncryptedPassword, u.Password) {
		return ErrInvalidEmailOrPassword
	} else {
		u.ID = found.ID
		u.Email = found.Email
		u.EncryptedPassword = found.EncryptedPassword
	}

	session, err := s.ss.Get(r, model.SessionName)
	if err != nil {
		return err
	}

	session.Values[model.SessionUserID] = u.ID

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func validate(u *model.User) error {
	if err := validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	); err != nil {
		return fmt.Errorf("%w: %w", ErrUserValidation, err)
	}
	return nil
}

func beforeSave(u *model.User) error {
	if len(u.Password) > 0 {
		encryptedPassword, err := encryptPassword(u.Password)
		if err != nil {
			return fmt.Errorf("before save: %w", err)
		}
		u.EncryptedPassword = encryptedPassword
	}
	return nil
}

func comparePassword(encryptedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)) == nil
}

func encryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("encrypt password: %w", err)
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
