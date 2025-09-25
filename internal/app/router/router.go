package router

import "net/http"

type Router interface {
	Handler() http.Handler
	Configure(createUser, createSession http.HandlerFunc)
}
