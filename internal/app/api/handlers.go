package api

import (
	"encoding/json"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/service/user"
)

// handleSignUp creates a new user
// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param input body UserCreateRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/signup [post]
func (a *api) handleSignUp() http.HandlerFunc {

	var service = user.New(a.db, a.ss, a.lg)

	return func(w http.ResponseWriter, r *http.Request) {

		req := new(UserCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := service.Save(r.Context(), u); err != nil {
			a.handleAllErrors(w, r, err)
		} else {
			a.respond(w, r, http.StatusCreated, &UserResponse{
				ID:    u.ID,
				Email: u.Email,
			})
		}
	}
}

// handleSignIn creates a session for the user
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
// @Router /users/signin [post]
func (a *api) handleSignIn() http.HandlerFunc {

	var service = user.New(a.db, a.ss, a.lg)

	return func(w http.ResponseWriter, r *http.Request) {
		req := new(SessionCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.errorNoLog(w, r, http.StatusBadRequest, err)
			return
		}

		if err := service.CreateHttpSession(w, r, &model.User{
			Email:    req.Email,
			Password: req.Password,
		}); err != nil {
			a.handleAllErrors(w, r, err)
		} else {
			a.respond(w, r, http.StatusOK, nil)
		}
	}
}

// handleMe returns information about the current authenticated user
// @Summary Information about the current authenticated user
// @Description Returns information about the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/me [get]
func (a *api) handleMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u, err := user.Me(r.Context()); err != nil {
			a.handleAllErrors(w, r, err)
		} else {
			a.respond(w, r, http.StatusOK, &UserResponse{
				ID:    u.ID,
				Email: u.Email,
			})
		}
	}
}
