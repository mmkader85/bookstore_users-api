package services

import (
	"fmt"

	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils"
)

var UserService userServiceInterface = &userServiceStruct{}

type userServiceStruct struct{}

type userServiceInterface interface {
	Delete(int64) *utils.RestErr
	Get(int64) (*users.User, *utils.RestErr)
	Search(string) ([]users.User, *utils.RestErr)
	Create(*users.User) (*users.User, *utils.RestErr)
	Update(bool, *users.User) (*users.User, *utils.RestErr)
}

func (userServiceStruct) Delete(userID int64) *utils.RestErr {
	var user users.User
	user.ID = userID
	if err := user.Delete(); err != nil {
		return err
	}

	return nil
}

func (userServiceStruct) Get(userID int64) (*users.User, *utils.RestErr) {
	var user users.User
	user.ID = userID
	if err := user.Get(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userServiceStruct) Search(status string) ([]users.User, *utils.RestErr) {
	user := &users.User{Status: status}

	return user.FindByStatus()
}

func (userServiceStruct) Create(user *users.User) (*users.User, *utils.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPwd, err := utils.PwdUtils.GeneratePwdHash(user.Password)
	if err != nil {
		return nil, utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Unable to generate password hash.%s", err.Error()))
	}
	user.Password = hashedPwd

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u userServiceStruct) Update(isPartial bool, user *users.User) (*users.User, *utils.RestErr) {
	currentUser, getUserErr := u.Get(user.ID)
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
