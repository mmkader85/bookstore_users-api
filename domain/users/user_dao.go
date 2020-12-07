package users

import (
	"database/sql"
	"fmt"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/utils/date"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
	"strings"
)

const emailUniqueKey	= "unq_email"
const insertUsersQuery	= "INSERT INTO users(first_name, last_name, email, created_at) VALUES(?, ?, ?, ?);"
const selectUserQuery	= "SELECT id, first_name, last_name, email, created_at FROM users WHERE id = ?;"

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(selectUserQuery)
	if err != nil {
		return errors.InternalServerErr("error preparing query for Get")
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.ID).Scan(&user.ID , &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return errors.NotFoundErr(fmt.Sprintf("userID %d doesn't exist", user.ID))
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	var (
		err     error
		stmt    *sql.Stmt
		result  sql.Result
	)

	stmt, err = users_db.Client.Prepare(insertUsersQuery)
	if err != nil {
		return errors.InternalServerErr("unable to prepare save query: " + err.Error())
	}
	defer stmt.Close()

	user.CreatedAt = date.GetNowString()
	result, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
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
