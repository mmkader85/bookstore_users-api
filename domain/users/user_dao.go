package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/logger"
	"github.com/mmkader85/bookstore_users-api/utils"
)

const (
	emailUniqueKey        = "unq_email"
	insertUsersQuery      = "INSERT INTO users(first_name, last_name, email, created_at, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	selectUserQuery       = "SELECT id, first_name, last_name, email, created_at, status, password FROM users WHERE id = ?;"
	updateUserQuery       = "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id = ?;"
	deleteUserQuery       = "DELETE FROM users WHERE id = ?;"
	findUserByStatusQuery = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status = ?;"
)

func (u *User) Save() *utils.RestErr {
	var (
		err          error
		stmt         *sql.Stmt
		result       sql.Result
		clientErrMsg = "Unable to save user"
	)

	stmt, err = users_db.Client.Prepare(insertUsersQuery)
	if err != nil {
		logger.Error("Unable to prepare save query.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}
	defer func() {
		_ = stmt.Close()
	}()

	u.CreatedAt = utils.DateUtils.GetNowDBFormat()
	result, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt, u.Status, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			logger.Info(fmt.Sprintf("Email '%s' already registered.", u.Email))
			return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Email '%s' already registered.", u.Email))
		}
		logger.Error("Unable to execute save query.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	u.ID, err = result.LastInsertId()
	if err != nil {
		logger.Error("Unable to get last insert id.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	return nil
}

func (u *User) Get() *utils.RestErr {
	clientErrMsg := "Unable to get user."
	stmt, err := users_db.Client.Prepare(selectUserQuery)
	if err != nil {
		logger.Error("Error preparing query for 'Get'.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.QueryRow(u.ID).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status, &u.Password)
	if err == sql.ErrNoRows {
		logger.Info(fmt.Sprintf("UserID %d doesn't exist", u.ID))
		return utils.RestErrUtils.NotFoundErr(fmt.Sprintf("UserID %d doesn't exist", u.ID))
	} else if err != nil {
		logger.Error("Error querying row for 'Get'.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	return nil
}

func (u *User) Update() *utils.RestErr {
	clientErrMsg := "Unable to update user."
	stmt, err := users_db.Client.Prepare(updateUserQuery)
	if err != nil {
		logger.Error("Error preparing query for 'Update'.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Status, u.ID)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			logger.Info(fmt.Sprintf("Email '%s' already exists.", u.Email))
			return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Email '%s' already exists.", u.Email))
		}
		logger.Error("Error executing query for 'Update'.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	return nil
}

func (u *User) Delete() *utils.RestErr {
	clientErrMsg := "Unable to delete user."
	stmt, err := users_db.Client.Prepare(deleteUserQuery)
	if err != nil {
		logger.Error("Error preparing query for 'Delete'.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}
	defer func() {
		_ = stmt.Close()
	}()

	res, deleteErr := stmt.Exec(u.ID)
	if deleteErr != nil {
		logger.Error("Error executing query for 'Delete'.", deleteErr)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	rowsAffected, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		logger.Error("Error getting RowsAffected() in 'Delete'.", rowsErr)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	} else if rowsAffected < 1 {
		logger.Info(fmt.Sprintf("UserID %d doesn't exist", u.ID))
		return utils.RestErrUtils.NotFoundErr(fmt.Sprintf("UserID %d doesn't exist", u.ID))
	}

	return nil
}

func (u *User) FindByStatus() ([]User, *utils.RestErr) {
	clientErrMsg := "Unable to search user by status."
	stmt, err := users_db.Client.Prepare(findUserByStatusQuery)
	if err != nil {
		logger.Error("Error preparing query for 'FindByStatus'.", err)
		return nil, utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, queryErr := stmt.Query(u.Status)
	if queryErr == sql.ErrNoRows {
		logger.Info(fmt.Sprintf("No user found with status %s.", u.Status))
		return nil, utils.RestErrUtils.NotFoundErr(fmt.Sprintf("No user found with status %s.", u.Status))
	} else if queryErr != nil {
		logger.Error("Error querying for 'FindByStatus'.", queryErr)
		return nil, utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}
	defer func() {
		_ = rows.Close()
	}()

	var users = make([]User, 0)
	for rows.Next() {
		var u User
		scanErr := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status)
		if scanErr != nil {
			logger.Error("Error scanning row in 'FindByStatus'.", scanErr)
			return nil, utils.RestErrUtils.InternalServerErr(clientErrMsg)
		}
		users = append(users, u)
	}

	return users, nil
}
