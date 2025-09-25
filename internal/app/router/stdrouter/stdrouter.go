package stdrouter

import "net/http"

type Router struct {
	h *http.ServeMux
}

func (r *Router) Handler() http.Handler {
	return r.h
}

func New() *Router {
	return &Router{h: http.NewServeMux()}
}

func (r *Router) Configure(createUser, createSession http.HandlerFunc) {
	r.h.Handle("/users", createUser)
	r.h.Handle("/sessions", createSession)
}
