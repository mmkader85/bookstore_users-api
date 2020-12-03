package users

import (
	"fmt"
	"github.com/mmkader85/bookstore_users-api/utils/date"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

var usersDB = make(map[int64]*User)

func (user *User) Get() *errors.RestErr {
	current := usersDB[user.Id]
	if current == nil {
		return errors.NotFoundErr(fmt.Sprintf("UserId %d doesn't exist", user.Id))
	}

	user.Id = current.Id
	user.FirstName = current.FirstName
	user.LastName = current.LastName
	user.Email = current.Email
	user.CreatedAt = current.CreatedAt

	return nil
}

func (user *User) Save() *errors.RestErr {
	current := usersDB[user.Id]
	if current != nil {
		return errors.BadRequestErr(fmt.Sprintf("UserId %d already exists", user.Id))
	}

	for _, existingUser := range usersDB {
		if user.Email == existingUser.Email {
			return errors.BadRequestErr(fmt.Sprintf("Email %s already registered", user.Email))
		}
	}

	dateNow := date.GetNowString()
	user.CreatedAt = dateNow
	usersDB[user.Id] = user

	return nil
}
