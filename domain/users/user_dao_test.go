package users

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	users_db.Client = db
	t.Cleanup(func() {
		_ = db.Close()
	})
	return mock
}

func TestUserSave_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("John", "Doe", "john@example.com", sqlmock.AnyArg(), "active", "hash").
		WillReturnResult(sqlmock.NewResult(7, 1))

	u := &User{FirstName: "John", LastName: "Doe", Email: "john@example.com", Status: "active", Password: "hash"}
	err := u.Save()

	assert.Nil(t, err)
	assert.EqualValues(t, 7, u.ID)
	assert.NotEmpty(t, u.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSave_PrepareError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("INSERT INTO users").WillReturnError(errors.New("boom"))

	err := (&User{}).Save()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, "Unable to save user", err.Message)
}

func TestUserSave_DuplicateEmail(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WillReturnError(errors.New("Error 1062: Duplicate entry 'a@b.c' for key 'unq_email'"))

	err := (&User{Email: "a@b.c"}).Save()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Contains(t, err.Message, "already registered")
}

func TestUserSave_ExecError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WillReturnError(errors.New("connection lost"))

	err := (&User{}).Save()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserSave_LastInsertIDError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(errors.New("no id")))

	err := (&User{}).Save()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserGet_Success(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "created_at", "status", "password"}).
		AddRow(1, "John", "Doe", "john@example.com", "2020-01-01 00:00:00", "active", "hash")
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WithArgs(int64(1)).WillReturnRows(rows)

	u := &User{ID: 1}
	err := u.Get()

	assert.Nil(t, err)
	assert.Equal(t, "John", u.FirstName)
	assert.Equal(t, "john@example.com", u.Email)
	assert.Equal(t, "hash", u.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserGet_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WithArgs(int64(9)).WillReturnError(sql.ErrNoRows)

	err := (&User{ID: 9}).Get()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status)
}

func TestUserGet_QueryError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WillReturnError(errors.New("boom"))

	err := (&User{ID: 1}).Get()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserGet_PrepareError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("SELECT id, first_name").WillReturnError(errors.New("boom"))

	err := (&User{ID: 1}).Get()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserUpdate_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("UPDATE users SET").
		ExpectExec().
		WithArgs("John", "Doe", "john@example.com", "active", int64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := (&User{ID: 3, FirstName: "John", LastName: "Doe", Email: "john@example.com", Status: "active"}).Update()

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserUpdate_DuplicateEmail(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("UPDATE users SET").
		ExpectExec().
		WillReturnError(errors.New("Duplicate entry for key 'unq_email'"))

	err := (&User{ID: 3, Email: "a@b.c"}).Update()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
}

func TestUserUpdate_ExecError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("UPDATE users SET").ExpectExec().WillReturnError(errors.New("boom"))

	err := (&User{ID: 3}).Update()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserUpdate_PrepareError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("UPDATE users SET").WillReturnError(errors.New("boom"))

	err := (&User{ID: 3}).Update()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserDelete_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("DELETE FROM users").ExpectExec().WithArgs(int64(5)).WillReturnResult(sqlmock.NewResult(0, 1))

	assert.Nil(t, (&User{ID: 5}).Delete())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserDelete_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("DELETE FROM users").ExpectExec().WillReturnResult(sqlmock.NewResult(0, 0))

	err := (&User{ID: 5}).Delete()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status)
}

func TestUserDelete_ExecError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("DELETE FROM users").ExpectExec().WillReturnError(errors.New("boom"))

	err := (&User{ID: 5}).Delete()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserDelete_RowsAffectedError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("DELETE FROM users").ExpectExec().WillReturnResult(sqlmock.NewErrorResult(errors.New("boom")))

	err := (&User{ID: 5}).Delete()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserDelete_PrepareError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("DELETE FROM users").WillReturnError(errors.New("boom"))

	err := (&User{ID: 5}).Delete()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserFindByStatus_Success(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "created_at", "status"}).
		AddRow(1, "A", "B", "a@b.c", "2020-01-01 00:00:00", "active").
		AddRow(2, "C", "D", "c@d.e", "2020-01-02 00:00:00", "active")
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WithArgs("active").WillReturnRows(rows)

	result, err := (&User{Status: "active"}).FindByStatus()

	assert.Nil(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "a@b.c", result[0].Email)
	assert.EqualValues(t, 2, result[1].ID)
}

func TestUserFindByStatus_Empty(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "created_at", "status"})
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WillReturnRows(rows)

	result, err := (&User{Status: "active"}).FindByStatus()

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestUserFindByStatus_PrepareError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("SELECT id, first_name").WillReturnError(errors.New("boom"))

	_, err := (&User{}).FindByStatus()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserFindByStatus_QueryError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WillReturnError(errors.New("boom"))

	_, err := (&User{}).FindByStatus()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestUserFindByStatus_ScanError(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "created_at", "status"}).
		AddRow("not-a-number", "A", "B", "a@b.c", "2020-01-01", "active")
	mock.ExpectPrepare("SELECT id, first_name").ExpectQuery().WillReturnRows(rows)

	_, err := (&User{}).FindByStatus()

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}
