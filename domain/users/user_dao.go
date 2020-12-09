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
	deleteUserQuery  = "DELETE FROM users WHERE id = ?;"
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
		return errors.InternalServerErr(fmt.Sprintf("Unable to prepare save query. \n%s", err.Error()))
	}

	user.CreatedAt = date_utils.GetNowString()
	result, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			return errors.BadRequestErr(fmt.Sprintf("Email '%s' already registered", user.Email))
		}
		return errors.InternalServerErr(fmt.Sprintf("Unable to execute save query. \n%s", err.Error()))
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return errors.InternalServerErr(fmt.Sprintf("Unable to get last insert id. \n%s", err.Error()))
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
		return errors.InternalServerErr(fmt.Sprintf("Error preparing query for 'Get' \n%s", err.Error()))
	}

	row := stmt.QueryRow(user.ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if row == sql.ErrNoRows {
		return errors.NotFoundErr(fmt.Sprintf("UserID %d doesn't exist", user.ID))
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
		return errors.InternalServerErr(fmt.Sprintf("Error preparing query for 'Update'. \n%s", err.Error()))
	}

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			return errors.BadRequestErr(fmt.Sprintf("Email '%s' already exists", user.Email))
		}
		return errors.InternalServerErr(fmt.Sprintf("Unable to execute query. \n %s", err.Error()))
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(deleteUserQuery)
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Println("Error closing statement:", err)
		}
	}()

	if err != nil {
		return errors.InternalServerErr(fmt.Sprintf("Error preparing query for 'Delete'. \n%s", err.Error()))
	}

	res, deleteErr := stmt.Exec(user.ID)
	if deleteErr != nil {
		return errors.InternalServerErr(fmt.Sprintf("Error deleting UserID: %d. \n%s", user.ID, deleteErr.Error()))
	}

	rowsAffected, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return errors.InternalServerErr(fmt.Sprintf("Unable to delete UserID %d. \n%s", user.ID, rowsErr.Error()))
	} else if rowsAffected < 1 {
		return errors.NotFoundErr(fmt.Sprintf("UserID %d doesn't exist", user.ID))
	}

	return nil
}
