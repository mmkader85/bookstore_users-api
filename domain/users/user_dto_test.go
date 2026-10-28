package users

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_Success(t *testing.T) {
	u := &User{FirstName: "  John ", LastName: " Doe  ", Email: "  John@Example.COM ", Status: " Active ", Password: "secret"}

	err := u.Validate()

	assert.Nil(t, err)
	assert.Equal(t, "John", u.FirstName)
	assert.Equal(t, "Doe", u.LastName)
	assert.Equal(t, "john@example.com", u.Email)
	assert.Equal(t, "active", u.Status)
}

func TestValidate_Errors(t *testing.T) {
	tests := []struct {
		name string
		user User
	}{
		{"empty email", User{Email: " ", Status: "active", Password: "x"}},
		{"invalid status", User{Email: "a@b.c", Status: "bogus", Password: "x"}},
		{"empty status", User{Email: "a@b.c", Password: "x"}},
		{"empty password", User{Email: "a@b.c", Status: "active", Password: "  "}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.user
			err := u.Validate()
			require.NotNil(t, err)
			assert.Equal(t, http.StatusBadRequest, err.Status)
		})
	}
}

func TestMarshal_Public(t *testing.T) {
	u := User{ID: 1, FirstName: "John", LastName: "Doe", Email: "a@b.c", CreatedAt: "t", Status: "active", Password: "hash"}

	out, ok := u.Marshal(true).(*PublicUser)

	require.True(t, ok)
	assert.Equal(t, &PublicUser{ID: 1, Email: "a@b.c", CreatedAt: "t", Status: "active"}, out)
}

func TestMarshal_Private(t *testing.T) {
	u := User{ID: 1, FirstName: "John", LastName: "Doe", Email: "a@b.c", CreatedAt: "t", Status: "active", Password: "hash"}

	out, ok := u.Marshal(false).(PrivateUser)

	require.True(t, ok)
	assert.Equal(t, PrivateUser{ID: 1, FirstName: "John", LastName: "Doe", Email: "a@b.c", CreatedAt: "t", Status: "active"}, out)
}
