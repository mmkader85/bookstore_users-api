package services

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var userCols = []string{"id", "first_name", "last_name", "email", "created_at", "status", "password"}

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

func TestGet_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WithArgs(int64(1), "deleted").
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(1, "John", "Doe", "j@e.c", "ts", "active", "hash"))

	user, err := UserService.Get(1)

	assert.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "j@e.c", user.Email)
}

func TestGet_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WillReturnRows(sqlmock.NewRows(userCols))

	user, err := UserService.Get(1)

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status)
}

func TestCreate_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectExec("INSERT INTO users").
		WithArgs("John", "Doe", "j@e.c", sqlmock.AnyArg(), "active", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(11, 1))

	in := &users.User{FirstName: "John", LastName: "Doe", Email: "J@E.c", Status: "active", Password: "Secret1"}
	user, err := UserService.Create(in)

	assert.Nil(t, err)
	require.NotNil(t, user)
	assert.EqualValues(t, 11, user.ID)
	assert.Equal(t, "j@e.c", user.Email)
	assert.NotEqual(t, "Secret1", user.Password, "password must be stored hashed")
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("Secret1")))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ValidationError(t *testing.T) {
	setupMockDB(t)

	user, err := UserService.Create(&users.User{Email: "", Status: "active", Password: "x"})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
}

type failingPwdUtils struct{}

func (failingPwdUtils) GeneratePwdHash(string) (string, error) { return "", errors.New("hash fail") }

func TestCreate_HashError(t *testing.T) {
	setupMockDB(t)
	orig := utils.PwdUtils
	utils.PwdUtils = failingPwdUtils{}
	defer func() { utils.PwdUtils = orig }()

	user, err := UserService.Create(&users.User{Email: "a@b.c", Status: "active", Password: "x"})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestCreate_SaveError(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectExec("INSERT INTO users").WillReturnError(errors.New("boom"))

	user, err := UserService.Create(&users.User{Email: "a@b.c", Status: "active", Password: "x"})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func expectGet(mock sqlmock.Sqlmock, id int64) {
	mock.ExpectQuery("SELECT id, first_name").WithArgs(id, "deleted").
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(id, "John", "Doe", "j@e.c", "ts", "active", "hash"))
}

func TestUpdate_Full(t *testing.T) {
	mock := setupMockDB(t)
	expectGet(mock, 1)
	mock.ExpectExec("UPDATE users SET").
		WithArgs("New", "Name", "n@e.c", "suspended", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user, err := UserService.Update(false, &users.User{ID: 1, FirstName: "New", LastName: "Name", Email: "n@e.c", Status: "suspended"})

	assert.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "New", user.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Partial(t *testing.T) {
	mock := setupMockDB(t)
	expectGet(mock, 1)
	mock.ExpectExec("UPDATE users SET").
		WithArgs("John", "Doe", "j@e.c", "suspended", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	user, err := UserService.Update(true, &users.User{ID: 1, Status: "suspended"})

	assert.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "suspended", user.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_PartialAllFields(t *testing.T) {
	mock := setupMockDB(t)
	expectGet(mock, 1)
	mock.ExpectExec("UPDATE users SET").
		WithArgs("A", "B", "c@d.e", "deleted", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := UserService.Update(true, &users.User{ID: 1, FirstName: "A", LastName: "B", Email: "c@d.e", Status: "deleted"})

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WillReturnRows(sqlmock.NewRows(userCols))

	user, err := UserService.Update(false, &users.User{ID: 1})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status)
}

func TestUpdate_ValidationError(t *testing.T) {
	mock := setupMockDB(t)
	expectGet(mock, 1)

	user, err := UserService.Update(false, &users.User{ID: 1, Email: "", Status: "active"})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
}

func TestUpdate_DBError(t *testing.T) {
	mock := setupMockDB(t)
	expectGet(mock, 1)
	mock.ExpectExec("UPDATE users SET").WillReturnError(errors.New("boom"))

	user, err := UserService.Update(false, &users.User{ID: 1, Email: "a@b.c", Status: "active"})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestDelete_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectExec("UPDATE users SET status").WithArgs("deleted", int64(4), "deleted").WillReturnResult(sqlmock.NewResult(0, 1))

	assert.Nil(t, UserService.Delete(4))
}

func TestDelete_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectExec("UPDATE users SET status").WillReturnResult(sqlmock.NewResult(0, 0))

	err := UserService.Delete(4)

	require.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status)
}

func TestSearch_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WithArgs("active", DefaultSearchLimit, 0).
		WillReturnRows(sqlmock.NewRows(userCols[:6]).AddRow(1, "A", "B", "a@b.c", "ts", "active"))

	result, err := UserService.Search("Active", 0, -3)

	assert.Nil(t, err)
	assert.Len(t, result, 1)
}

func TestSearch_Error(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WillReturnError(errors.New("boom"))

	_, err := UserService.Search("active", 0, 0)

	require.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}

func TestCreate_NilUser(t *testing.T) {
	setupMockDB(t)

	user, err := UserService.Create(nil)

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
}

func TestCreate_EmptyPassword(t *testing.T) {
	setupMockDB(t)

	user, err := UserService.Create(&users.User{Email: "a@b.c", Status: "active", Password: "  "})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
}

func TestCreate_PasswordTooLongForBcrypt(t *testing.T) {
	setupMockDB(t)

	user, err := UserService.Create(&users.User{Email: "a@b.c", Status: "active", Password: strings.Repeat("a", 100)})

	assert.Nil(t, user)
	require.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
}

func TestUpdate_DoesNotRequirePasswordOrRehash(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WithArgs(int64(1), "deleted").
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(1, "John", "Doe", "j@e.c", "ts", "active", ""))
	mock.ExpectExec("UPDATE users SET").WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := UserService.Update(true, &users.User{ID: 1, FirstName: "X"})

	assert.Nil(t, err)
}

func TestSearch_InvalidStatus(t *testing.T) {
	setupMockDB(t)

	for _, status := range []string{"", "bogus"} {
		result, err := UserService.Search(status, 10, 0)
		assert.Nil(t, result)
		require.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Status)
	}
}

func TestSearch_LimitIsCapped(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery("SELECT id, first_name").WithArgs("active", MaxSearchLimit, 40).
		WillReturnRows(sqlmock.NewRows(userCols[:6]))

	_, err := UserService.Search("active", 1000000, 40)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
