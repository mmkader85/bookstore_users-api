package users

import (
	"database/sql"
	"fmt"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/utils/date"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
	"strings"
)

var usersDB = make(map[int64]*User)
const insertUsersQuery = "INSERT INTO users(first_name, last_name, email, created_at) VALUES(?, ?, ?, ?);"

func (user *User) Get() *errors.RestErr {
	current := usersDB[user.ID]
	if current == nil {
		return errors.NotFoundErr(fmt.Sprintf("userID %d doesn't exist", user.ID))
	}

	user.ID = current.ID
	user.FirstName = current.FirstName
	user.LastName = current.LastName
	user.Email = current.Email
	user.CreatedAt = current.CreatedAt

	return nil
}

func (user *User) Save() *errors.RestErr {
	var (
		err error
		stmt *sql.Stmt
		result sql.Result
	)

	stmt, err = users_db.Client.Prepare(insertUsersQuery)
	if err != nil {
		return errors.InternalServerErr("unable to prepare save query: " + err.Error())
	}
	defer stmt.Close()

	user.CreatedAt = date.GetNowString()
	result, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unq_email") {
			return errors.BadRequestErr(fmt.Sprintf("email '%s' already registered", user.Email))
		}
		return errors.InternalServerErr("unable to execute query: " + err.Error())
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return errors.InternalServerErr("unable to get last insert id: " + err.Error())
	}

	return nil
}
