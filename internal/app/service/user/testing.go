package user

import (
	"log/slog"
	"os"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/storage/memorystorage"
	"github.com/gorilla/sessions"
)

var TestSessionKey = []byte("session_key")

func TestNew(t *testing.T) *UserService {
	t.Helper()

	var (
		userStorage  = memorystorage.New()
		sessionStore = sessions.NewCookieStore(TestSessionKey)
		logger       = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	return New(userStorage, sessionStore, logger)
}
