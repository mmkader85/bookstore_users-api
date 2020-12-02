package users

import (
	"github.com/mmkader85/bookstore_users-api/utils/errors"
	"strings"
)

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required"`
	CreatedAt string `json:"created_at"`
}

func (user *User) Validate() *errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.BadRequestErr("empty user email")
	}

	return nil
}
