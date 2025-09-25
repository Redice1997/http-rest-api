package gorillarouter

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

func (r *Router) Handler() http.Handler {
	return r.router
}

func (r *Router) Configure(createUser, createSession http.HandlerFunc) {
	r.router.HandleFunc("/users", createUser).Methods("POST")
	r.router.HandleFunc("/sessions", createSession).Methods("POST")
}
