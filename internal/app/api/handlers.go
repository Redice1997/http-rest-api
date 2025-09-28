package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

// handleUserCreate creates a new user
// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param input body UserCreateRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/users [post]
func (a *api) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(UserCreateRequest)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := a.db.User().GetByEmail(r.Context(), req.Email); errors.Is(err, storage.ErrRecordNotFound) {
		} else if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		} else {
			a.error(w, r, http.StatusBadRequest, ErrAlreadyExists)
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

		a.respond(w, r, http.StatusCreated, &UserResponse{
			ID:    u.ID,
			Email: u.Email,
		})
	}
}

// handleSessionCreate creates a session for the user
// @Summary Enter the system
// @Description Creates a session for the user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body SessionCreateRequest true "Pass data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/sessions [post]
func (a *api) handleSessionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(SessionCreateRequest)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.errorNoLog(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := a.db.User().GetByEmail(r.Context(), req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			a.errorNoLog(w, r, http.StatusUnauthorized, ErrInvalidEmailOrPassword)
			return
		}

		session, err := a.sess.Get(r, SessionName)
		if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID

		if err := session.Save(r, w); err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.respond(w, r, http.StatusOK, nil)
	}
}

// handleWhoAmI returns information about the current authenticated user
// @Summary Information about the current authenticated user
// @Description Returns information about the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/whoami [get]
func (a *api) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(CtxUserKey).(*model.User)
		if u == nil {
			a.errorNoLog(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		a.respond(w, r, http.StatusOK, &UserResponse{
			ID:    u.ID,
			Email: u.Email,
		})
	}
}
