package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/services"
	"github.com/mmkader85/bookstore_users-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	deleteFn func(int64) *utils.RestErr
	getFn    func(int64) (*users.User, *utils.RestErr)
	searchFn func(string, int, int) ([]users.User, *utils.RestErr)
	createFn func(*users.User) (*users.User, *utils.RestErr)
	updateFn func(bool, *users.User) (*users.User, *utils.RestErr)
}

func (m mockService) Delete(id int64) *utils.RestErr             { return m.deleteFn(id) }
func (m mockService) Get(id int64) (*users.User, *utils.RestErr) { return m.getFn(id) }
func (m mockService) Search(s string, l, o int) ([]users.User, *utils.RestErr) {
	return m.searchFn(s, l, o)
}
func (m mockService) Create(u *users.User) (*users.User, *utils.RestErr) { return m.createFn(u) }
func (m mockService) Update(p bool, u *users.User) (*users.User, *utils.RestErr) {
	return m.updateFn(p, u)
}

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/user", Create)
	r.GET("/user/:user_id", Get)
	r.PUT("/user/:user_id", Update)
	r.PATCH("/user/:user_id", Update)
	r.DELETE("/user/:user_id", Delete)
	r.GET("/internal/user/search", Search)
	return r
}

func withService(t *testing.T, m mockService) {
	t.Helper()
	orig := services.UserService
	services.UserService = m
	t.Cleanup(func() { services.UserService = orig })
}

func do(r http.Handler, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func sampleUser() *users.User {
	return &users.User{ID: 1, FirstName: "John", LastName: "Doe", Email: "j@e.c", CreatedAt: "ts", Status: "active", Password: "hash"}
}

func TestCreate_Success_Public(t *testing.T) {
	withService(t, mockService{createFn: func(u *users.User) (*users.User, *utils.RestErr) {
		assert.Equal(t, "j@e.c", u.Email)
		return sampleUser(), nil
	}})

	w := do(newRouter(), http.MethodPost, "/user", `{"email":"j@e.c","status":"active","password":"x"}`, nil)

	assert.Equal(t, http.StatusCreated, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "j@e.c", body["email"])
	assert.NotContains(t, body, "first_name")
	assert.NotContains(t, body, "password")
}

func TestCreate_Success_Private(t *testing.T) {
	withService(t, mockService{createFn: func(u *users.User) (*users.User, *utils.RestErr) { return sampleUser(), nil }})

	w := do(newRouter(), http.MethodPost, "/user", `{"email":"j@e.c"}`, map[string]string{"X-Private": "true"})

	assert.Equal(t, http.StatusCreated, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "John", body["first_name"])
	assert.NotContains(t, body, "password")
}

func TestCreate_InvalidJSON(t *testing.T) {
	withService(t, mockService{})

	w := do(newRouter(), http.MethodPost, "/user", `{bad`, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_MissingEmail_RejectedByService(t *testing.T) {
	withService(t, mockService{createFn: func(*users.User) (*users.User, *utils.RestErr) {
		return nil, utils.RestErrUtils.BadRequestErr("Empty user email.")
	}})

	w := do(newRouter(), http.MethodPost, "/user", `{"first_name":"x"}`, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_ServiceError(t *testing.T) {
	withService(t, mockService{createFn: func(*users.User) (*users.User, *utils.RestErr) {
		return nil, utils.RestErrUtils.BadRequestErr("nope")
	}})

	w := do(newRouter(), http.MethodPost, "/user", `{"email":"j@e.c"}`, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "nope")
}

func TestGet_Success(t *testing.T) {
	withService(t, mockService{getFn: func(id int64) (*users.User, *utils.RestErr) {
		assert.EqualValues(t, 1, id)
		return sampleUser(), nil
	}})

	w := do(newRouter(), http.MethodGet, "/user/1", "", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"email":"j@e.c"`)
}

func TestGet_InvalidID(t *testing.T) {
	withService(t, mockService{})

	w := do(newRouter(), http.MethodGet, "/user/abc", "", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGet_NotFound(t *testing.T) {
	withService(t, mockService{getFn: func(int64) (*users.User, *utils.RestErr) {
		return nil, utils.RestErrUtils.NotFoundErr("missing")
	}})

	w := do(newRouter(), http.MethodGet, "/user/1", "", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdate_PutAndPatch(t *testing.T) {
	for method, partial := range map[string]bool{http.MethodPut: false, http.MethodPatch: true} {
		method, partial := method, partial
		t.Run(method, func(t *testing.T) {
			withService(t, mockService{updateFn: func(p bool, u *users.User) (*users.User, *utils.RestErr) {
				assert.Equal(t, partial, p)
				assert.EqualValues(t, 5, u.ID)
				return sampleUser(), nil
			}})

			w := do(newRouter(), method, "/user/5", `{"email":"j@e.c"}`, nil)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestUpdate_InvalidID(t *testing.T) {
	withService(t, mockService{})

	w := do(newRouter(), http.MethodPut, "/user/abc", `{"email":"j@e.c"}`, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	withService(t, mockService{})

	w := do(newRouter(), http.MethodPut, "/user/1", `{bad`, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	withService(t, mockService{updateFn: func(bool, *users.User) (*users.User, *utils.RestErr) {
		return nil, utils.RestErrUtils.InternalServerErr("boom")
	}})

	w := do(newRouter(), http.MethodPut, "/user/1", `{"email":"j@e.c"}`, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDelete_Success(t *testing.T) {
	withService(t, mockService{deleteFn: func(id int64) *utils.RestErr {
		assert.EqualValues(t, 3, id)
		return nil
	}})

	w := do(newRouter(), http.MethodDelete, "/user/3", "", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "successfully deleted")
}

func TestDelete_InvalidID(t *testing.T) {
	withService(t, mockService{})

	assert.Equal(t, http.StatusBadRequest, do(newRouter(), http.MethodDelete, "/user/x", "", nil).Code)
}

func TestDelete_ServiceError(t *testing.T) {
	withService(t, mockService{deleteFn: func(int64) *utils.RestErr { return utils.RestErrUtils.NotFoundErr("missing") }})

	assert.Equal(t, http.StatusNotFound, do(newRouter(), http.MethodDelete, "/user/3", "", nil).Code)
}

func TestSearch_Success(t *testing.T) {
	withService(t, mockService{searchFn: func(s string, _, _ int) ([]users.User, *utils.RestErr) {
		assert.Equal(t, "active", s)
		return []users.User{*sampleUser(), *sampleUser()}, nil
	}})

	w := do(newRouter(), http.MethodGet, "/internal/user/search?status=active", "", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var body []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Len(t, body, 2)
}

func TestSearch_Empty(t *testing.T) {
	withService(t, mockService{searchFn: func(string, int, int) ([]users.User, *utils.RestErr) { return []users.User{}, nil }})

	w := do(newRouter(), http.MethodGet, "/internal/user/search?status=active", "", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[]`, w.Body.String())
}

func TestSearch_Error(t *testing.T) {
	withService(t, mockService{searchFn: func(string, int, int) ([]users.User, *utils.RestErr) {
		return nil, utils.RestErrUtils.InternalServerErr("boom")
	}})

	assert.Equal(t, http.StatusInternalServerError, do(newRouter(), http.MethodGet, "/internal/user/search", "", nil).Code)
}

func TestCreate_NullBody_NoPanic(t *testing.T) {
	withService(t, mockService{createFn: func(u *users.User) (*users.User, *utils.RestErr) {
		require.NotNil(t, u)
		return nil, utils.RestErrUtils.BadRequestErr("Empty user email.")
	}})

	w := do(newRouter(), http.MethodPost, "/user", `null`, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_BodyIDCannotOverrideURLID(t *testing.T) {
	withService(t, mockService{updateFn: func(_ bool, u *users.User) (*users.User, *utils.RestErr) {
		assert.EqualValues(t, 5, u.ID, "URL id must win over body id")
		return sampleUser(), nil
	}})

	w := do(newRouter(), http.MethodPut, "/user/5", `{"id":999,"email":"j@e.c"}`, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate_PatchWithoutEmail(t *testing.T) {
	withService(t, mockService{updateFn: func(p bool, u *users.User) (*users.User, *utils.RestErr) {
		assert.True(t, p)
		assert.Equal(t, "Jane", u.FirstName)
		return sampleUser(), nil
	}})

	w := do(newRouter(), http.MethodPatch, "/user/5", `{"first_name":"Jane"}`, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearch_ParsesPagination(t *testing.T) {
	withService(t, mockService{searchFn: func(s string, limit, offset int) ([]users.User, *utils.RestErr) {
		assert.Equal(t, 20, limit)
		assert.Equal(t, 40, offset)
		return []users.User{}, nil
	}})

	w := do(newRouter(), http.MethodGet, "/internal/user/search?status=active&limit=20&offset=40", "", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearch_IgnoresGarbagePagination(t *testing.T) {
	withService(t, mockService{searchFn: func(s string, limit, offset int) ([]users.User, *utils.RestErr) {
		assert.Zero(t, limit)
		assert.Zero(t, offset)
		return []users.User{}, nil
	}})

	w := do(newRouter(), http.MethodGet, "/internal/user/search?status=active&limit=abc&offset=xyz", "", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}
