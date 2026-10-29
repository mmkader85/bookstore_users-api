package users

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/mmkader85/bookstore_users-api/utils"
)

const (
	UserStatusActive    = "active"
	UserStatusDeleted   = "deleted"
	UserStatusSuspended = "suspended"

	// bcrypt silently/explicitly rejects input beyond 72 bytes.
	maxPasswordBytes = 72
	maxEmailLen      = 254
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Password  string `json:"password"`
}

type PublicUser struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type PrivateUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

func (u *User) Validate() *utils.RestErr {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return utils.RestErrUtils.BadRequestErr("Empty user email.")
	}
	if parsed, err := mail.ParseAddress(u.Email); err != nil || parsed.Address != u.Email || len(u.Email) > maxEmailLen {
		return utils.RestErrUtils.BadRequestErr("Invalid user email.")
	}

	u.Status = strings.TrimSpace(strings.ToLower(u.Status))
	if u.Status != UserStatusActive && u.Status != UserStatusSuspended && u.Status != UserStatusDeleted {
		return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("User status should be %s/%s/%s only.", UserStatusActive, UserStatusSuspended, UserStatusDeleted))
	}

	return nil
}

// ValidatePassword checks a plain-text password. The password is deliberately
// not trimmed or case-folded, as that would silently weaken what the user chose.
func (u *User) ValidatePassword() *utils.RestErr {
	if strings.TrimSpace(u.Password) == "" {
		return utils.RestErrUtils.BadRequestErr("Empty user password.")
	}
	if len(u.Password) > maxPasswordBytes {
		return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("User password should not exceed %d bytes.", maxPasswordBytes))
	}

	return nil
}

func (u User) Marshal(isPublic bool) interface{} {
	if isPublic {
		return &PublicUser{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			Status:    u.Status,
		}
	}

	return PrivateUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		Status:    u.Status,
	}
}
