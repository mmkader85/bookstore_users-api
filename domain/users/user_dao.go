package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
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
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)

	stmt, err = users_db.Client.Prepare(insertUsersQuery)
	if err != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Unable to prepare save query. %s", err.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	u.CreatedAt = utils.DateUtils.GetNowDBFormat()
	result, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.CreatedAt, u.Status, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Email '%s' already registered", u.Email))
		}
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Unable to execute save query. %s", err.Error()))
	}

	u.ID, err = result.LastInsertId()
	if err != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Unable to get last insert id. %s", err.Error()))
	}

	return nil
}

func (u *User) Get() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(selectUserQuery)
	if err != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error preparing query for 'Get' %s", err.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.QueryRow(u.ID).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status, &u.Password)
	if err == sql.ErrNoRows {
		return utils.RestErrUtils.NotFoundErr(fmt.Sprintf("UserID %d doesn't exist", u.ID))
	} else if err != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error querying row for 'Get' %s", err.Error()))
	}

	return nil
}

func (u *User) Update() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(updateUserQuery)
	if err != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error preparing query for 'Update'. %s", err.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Status, u.ID)
	if err != nil {
		if strings.Contains(err.Error(), emailUniqueKey) {
			return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Email '%s' already exists", u.Email))
		}
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Unable to execute query. %s", err.Error()))
	}

	return nil
}

func (u *User) Delete() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(deleteUserQuery)
	if err != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error preparing query for 'Delete'. %s", err.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	res, deleteErr := stmt.Exec(u.ID)
	if deleteErr != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error deleting UserID: %d. %s", u.ID, deleteErr.Error()))
	}

	rowsAffected, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Unable to delete UserID %d. %s", u.ID, rowsErr.Error()))
	} else if rowsAffected < 1 {
		return utils.RestErrUtils.NotFoundErr(fmt.Sprintf("UserID %d doesn't exist", u.ID))
	}

	return nil
}

func (u *User) FindByStatus() ([]User, *utils.RestErr) {
	stmt, err := users_db.Client.Prepare(findUserByStatusQuery)
	if err != nil {
		return nil, utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error preparing query for 'FindByStatus'. %s", err.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, queryErr := stmt.Query(u.Status)
	if queryErr == sql.ErrNoRows {
		return nil, utils.RestErrUtils.NotFoundErr(fmt.Sprintf("No user found with status %s.", u.Status))
	} else if queryErr != nil {
		return nil, utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error querying for 'FindByStatus'. %s", queryErr.Error()))
	}
	defer func() {
		_ = rows.Close()
	}()

	var users = make([]User, 0)
	for rows.Next() {
		var u User
		scanErr := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status)
		if scanErr != nil {
			return nil, utils.RestErrUtils.InternalServerErr(fmt.Sprintf("Error scanning row in 'FindByStatus'. %s", scanErr.Error()))
		}
		users = append(users, u)
	}

	return users, nil
}
