package api

import (
	"encoding/json"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/service/user"
)

// HandleSignUp creates a new user
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
func (a *API) HandleSignUp() http.HandlerFunc {
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
		if err := a.user.Save(r.Context(), u); err != nil {
			a.handleAllErrors(w, r, err)
		} else {
			a.respond(w, r, http.StatusCreated, &UserResponse{
				ID:    u.ID,
				Email: u.Email,
			})
		}
	}
}

// HandleSignIn creates a session for the user
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
func (a *API) HandleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(SessionCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.errorNoLog(w, r, http.StatusBadRequest, err)
			return
		}

		if err := a.user.CreateHttpSession(w, r, &model.User{
			Email:    req.Email,
			Password: req.Password,
		}); err != nil {
			a.handleAllErrors(w, r, err)
		} else {
			a.respond(w, r, http.StatusOK, nil)
		}
	}
}

// HandleMe returns information about the current authenticated user
// @Summary Information about the current authenticated user
// @Description Returns information about the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/me [get]
func (a *API) HandleMe() http.HandlerFunc {
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
