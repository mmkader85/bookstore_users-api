package services

import (
	"fmt"

	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
	"github.com/mmkader85/bookstore_users-api/utils/password_utils"
)

func CreateUser(user *users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPwd, err := password_utils.GeneratePwdHash(user.Password)
	if err != nil {
		return nil, errors.InternalServerErr(fmt.Sprintf("Unable to generate password hash.%s", err.Error()))
	}
	user.Password = hashedPwd

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
		if user.Status != "" {
			currentUser.Status = user.Status
		}
	} else {
		currentUser.FirstName = user.FirstName
		currentUser.LastName = user.LastName
		currentUser.Email = user.Email
		currentUser.Status = user.Status
	}

	if err := currentUser.Validate(); err != nil {
		return nil, err
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

func Search(status string) ([]map[string]string, *errors.RestErr) {
	user := &users.User{Status: status}

	return user.FindByStatus()
}
