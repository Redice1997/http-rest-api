package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/model"
)

var (
	errInvalidEmailOrPassword = errors.New("invalid email or password")
)

func (a *api) handleUserCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			a.error(w, r, http.StatusMethodNotAllowed, nil)
			return
		}

		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := a.db.User().Create(r.Context(), u); err != nil {
			a.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		a.respond(w, r, http.StatusCreated, &response{
			ID:    u.ID,
			Email: u.Email,
		})
	}
}

func (a *api) handleSessionCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			a.error(w, r, http.StatusMethodNotAllowed, nil)
			return
		}

		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := a.db.User().GetByEmail(r.Context(), req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			a.error(w, r, http.StatusUnauthorized, errInvalidEmailOrPassword)
			return
		}

		a.respond(w, r, http.StatusOK, nil)
	}
}
