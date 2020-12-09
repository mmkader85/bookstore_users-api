package services

import (
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

func CreateUser(user *users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUser(userID int64) (*users.User, *errors.RestErr) {
	var user users.User
	user.ID = userID
	if err := user.Get(); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(isPartial bool, user *users.User) (*users.User, *errors.RestErr) {
	currentUser, getUserErr := GetUser(user.ID)
	if getUserErr != nil {
		return currentUser, getUserErr
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			currentUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currentUser.LastName = user.LastName
		}
		if user.Email != "" {
			currentUser.Email = user.Email
		}
	} else {
		currentUser.FirstName = user.FirstName
		currentUser.LastName = user.LastName
		currentUser.Email = user.Email
	}

	if err := currentUser.Update(); err != nil {
		return nil, err
	}

	return currentUser, nil
}

func DeleteUser(userID int64) *errors.RestErr {
	var user users.User
	user.ID = userID
	if err := user.Delete(); err != nil {
		return err
	}

	return nil
}
