package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/logger"
	"github.com/mmkader85/bookstore_users-api/utils"
)

const (
	emailUniqueKey        = "unq_email"
	mysqlDuplicateEntry   = 1062
	insertUsersQuery      = "INSERT INTO users(first_name, last_name, email, created_at, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	selectUserQuery       = "SELECT id, first_name, last_name, email, created_at, status, password FROM users WHERE id = ? AND status <> ?;"
	updateUserQuery       = "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id = ?;"
	deleteUserQuery       = "UPDATE users SET status = ? WHERE id = ? AND status <> ?;"
	findUserByStatusQuery = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status = ? ORDER BY id LIMIT ? OFFSET ?;"
)

func isDuplicateEmail(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == mysqlDuplicateEntry
	}

	return strings.Contains(err.Error(), emailUniqueKey)
}

func (u *User) Save() *utils.RestErr {
	var (
		err          error
		result       sql.Result
		clientErrMsg = "Unable to save user"
	)

	u.CreatedAt = utils.DateUtils.GetNowDBFormat()
	result, err = users_db.Client.Exec(insertUsersQuery, u.FirstName, u.LastName, u.Email, u.CreatedAt, u.Status, u.Password)
	if err != nil {
		if isDuplicateEmail(err) {
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
	err := users_db.Client.QueryRow(selectUserQuery, u.ID, UserStatusDeleted).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.Status, &u.Password)
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
	_, err := users_db.Client.Exec(updateUserQuery, u.FirstName, u.LastName, u.Email, u.Status, u.ID)
	if err != nil {
		if isDuplicateEmail(err) {
			logger.Info(fmt.Sprintf("Email '%s' already exists.", u.Email))
			return utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Email '%s' already exists.", u.Email))
		}
		logger.Error("Error executing query for 'Update'.", err)
		return utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	return nil
}

// Delete soft-deletes the user by marking it as deleted. Already deleted users are treated as not found.
func (u *User) Delete() *utils.RestErr {
	clientErrMsg := "Unable to delete user."
	res, deleteErr := users_db.Client.Exec(deleteUserQuery, UserStatusDeleted, u.ID, UserStatusDeleted)
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

func (u *User) FindByStatus(limit, offset int) ([]User, *utils.RestErr) {
	clientErrMsg := "Unable to search user by status."
	rows, queryErr := users_db.Client.Query(findUserByStatusQuery, u.Status, limit, offset)
	if queryErr != nil {
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
	if rowsErr := rows.Err(); rowsErr != nil {
		logger.Error("Error iterating rows in 'FindByStatus'.", rowsErr)
		return nil, utils.RestErrUtils.InternalServerErr(clientErrMsg)
	}

	return users, nil
}
