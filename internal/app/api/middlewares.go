package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

func (a *api) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.sess.Get(r, SessionName)
		if err != nil {
			a.error(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		}
		id, ok := session.Values["user_id"].(int64)
		if !ok {
			a.error(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		ctx := r.Context()

		u, err := a.db.User().GetByID(ctx, id)
		if errors.Is(err, storage.ErrRecordNotFound) {
			a.error(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		} else if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, CtxUserKey, u)))
	})
}
