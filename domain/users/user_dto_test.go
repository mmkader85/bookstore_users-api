package users

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_Success(t *testing.T) {
	u := &User{FirstName: "  John ", LastName: " Doe  ", Email: "  John@Example.COM ", Status: " Active "}

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
		{"empty email", User{Email: " ", Status: "active"}},
		{"malformed email", User{Email: "not-an-email", Status: "active"}},
		{"email with display name", User{Email: "Bob <bob@example.com>", Status: "active"}},
		{"overlong email", User{Email: strings.Repeat("a", 250) + "@b.cc", Status: "active"}},
		{"invalid status", User{Email: "a@b.c", Status: "bogus"}},
		{"empty status", User{Email: "a@b.c"}},
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

func TestValidate_DoesNotTouchPassword(t *testing.T) {
	u := &User{Email: "a@b.c", Status: "active", Password: "  MixedCase Pass  "}

	assert.Nil(t, u.Validate())
	assert.Equal(t, "  MixedCase Pass  ", u.Password)
}

func TestValidatePassword(t *testing.T) {
	assert.Nil(t, (&User{Password: "MixedCase1"}).ValidatePassword())
	assert.Nil(t, (&User{Password: strings.Repeat("a", 72)}).ValidatePassword())

	for _, pwd := range []string{"", "   ", strings.Repeat("a", 73)} {
		err := (&User{Password: pwd}).ValidatePassword()
		require.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Status)
	}
}
