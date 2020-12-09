package users

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/utils/date_utils"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

const (
	emailUniqueKey   = "unq_email"
	insertUsersQuery = "INSERT INTO users(first_name, last_name, email, created_at) VALUES(?, ?, ?, ?);"
	selectUserQuery  = "SELECT id, first_name, last_name, email, created_at FROM users WHERE id = ?;"
	updateUserQuery  = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
)

func (user *User) Save() *errors.RestErr {
	var (
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)

	stmt, err = users_db.Client.Prepare(insertUsersQuery)
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Println("Error closing statement:", err)
		}
	}()
	if err != nil {
		return errors.InternalServerErr("Unable to prepare save query: " + err.Error())
	}

	user.CreatedAt = date_utils.GetNowString()
	result, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			return errors.BadRequestErr(fmt.Sprintf("Email '%s' already registered", user.Email))
		}
		return errors.InternalServerErr("Unable to execute query: " + err.Error())
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return errors.InternalServerErr("Unable to get last insert id: " + err.Error())
	}

	return nil
}

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(selectUserQuery)
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Println("Error closing statement:", err)
		}
	}()
	if err != nil {
		return errors.InternalServerErr("error preparing query for Get")
	}

	err = stmt.QueryRow(user.ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return errors.NotFoundErr(fmt.Sprintf("userID %d doesn't exist", user.ID))
	}

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(updateUserQuery)
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Println("Error closing statement:", err)
		}
	}()
	if err != nil {
		return errors.InternalServerErr("Error preparing query for Update")
	}

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			return errors.BadRequestErr(fmt.Sprintf("Email '%s' already exists", user.Email))
		}
		return errors.InternalServerErr("Unable to execute query: " + err.Error())
	}

	return nil
}
