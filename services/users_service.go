package services

import (
	"fmt"
	"strings"

	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils"
)

const (
	DefaultSearchLimit = 100
	MaxSearchLimit     = 1000
)

var UserService userServiceInterface = &userServiceStruct{}

type userServiceStruct struct{}

type userServiceInterface interface {
	Delete(int64) *utils.RestErr
	Get(int64) (*users.User, *utils.RestErr)
	Search(status string, limit, offset int) ([]users.User, *utils.RestErr)
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

func (userServiceStruct) Search(status string, limit, offset int) ([]users.User, *utils.RestErr) {
	status = strings.TrimSpace(strings.ToLower(status))
	if status != users.UserStatusActive && status != users.UserStatusSuspended && status != users.UserStatusDeleted {
		return nil, utils.RestErrUtils.BadRequestErr(fmt.Sprintf("User status should be %s/%s/%s only.", users.UserStatusActive, users.UserStatusSuspended, users.UserStatusDeleted))
	}
	if limit <= 0 {
		limit = DefaultSearchLimit
	} else if limit > MaxSearchLimit {
		limit = MaxSearchLimit
	}
	if offset < 0 {
		offset = 0
	}

	user := &users.User{Status: status}

	return user.FindByStatus(limit, offset)
}

func (userServiceStruct) Create(user *users.User) (*users.User, *utils.RestErr) {
	if user == nil {
		return nil, utils.RestErrUtils.BadRequestErr("Empty request body.")
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := user.ValidatePassword(); err != nil {
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
		return nil, getUserErr
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
