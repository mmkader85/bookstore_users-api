package services

import (
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

func CreateUser(user users.Users) (*users.Users, *errors.RestErr) {
	return &user, nil
}