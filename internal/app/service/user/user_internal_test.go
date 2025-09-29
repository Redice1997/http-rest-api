package user

import (
	"strings"
	"testing"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeSave(t *testing.T) {
	user := model.TestUser(t)
	assert.NoError(t, beforeSave(user))
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestUser_Validate(t *testing.T) {

	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid user",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "encrypted_password"
				return u
			},
			isValid: true,
		},
		{
			name: "empty password user",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email user",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email format",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "abc.org"
				return u
			},
			isValid: false,
		},
		{
			name: "short password user",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "pas"
				return u
			},
			isValid: false,
		},
		{
			name: "long password user",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = strings.Repeat("a", 101)
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, validate(tc.u()))
			} else {
				assert.Error(t, validate(tc.u()))
			}
		})
	}
}
