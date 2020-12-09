package services

import (
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

func GetUser(userID int64) (*users.User, *errors.RestErr) {
	var user users.User
	user.ID = userID
	if err := user.Get(); err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(user *users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Update(); err != nil {
		return nil, err
	}

	return user, nil
}
