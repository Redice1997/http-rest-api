package api

import "net/http"

func (a *api) handleHello() http.HandlerFunc {

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		a.response(w, http.StatusOK, &response{Message: "Hello, World!"})
	}
}
