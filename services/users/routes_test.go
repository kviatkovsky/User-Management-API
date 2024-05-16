package users

import (
	"encoding/json"
	"main/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("Should return status code 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/getlist", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/getlist", handler.handleGetList).Methods(http.MethodGet)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Should return status code 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/getlist", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/getlist", handler.handleGetList).Methods(http.MethodGet)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	})
}

func TestCreateUser(t *testing.T) {
	payload := types.UserPayload{
		User: types.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "password",
		},
		UserRole: &types.UserRole{
			AccessLevel: 3,
		},
	}

	bPayload, _ := json.Marshal(payload)

	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("handleCreate User, should return status code 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(string(bPayload)))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/create", handler.handleCreate).Methods(http.MethodPost)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("handleCreate User, Should return status code 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/create", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/create", handler.handleCreate).Methods(http.MethodPost)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	})
}

func TestEditUser(t *testing.T) {
	payload := types.UserPayload{
		User: types.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "password",
		},
		UserRole: &types.UserRole{
			AccessLevel: 3,
		},
	}

	bPayload, _ := json.Marshal(payload)

	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("handleEdit User, should return status code 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/edit/22", strings.NewReader(string(bPayload)))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/edit/{id}", handler.handleEdit).Methods(http.MethodPost)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("handleEdit User, Should return status code 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/edit/22", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/user/edit/{id}", handler.handleEdit).Methods(http.MethodPost)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) DeleteUserById(id string) {
}

func (m *mockUserStore) GetUserRoleByUserId(uid uint) *types.UserRole {
	return &types.UserRole{}
}

func (m *mockUserStore) UpdateUserAccessLevel(u *types.User, accessLevel uint) {
}

func (m *mockUserStore) GetList() []types.User {
	return nil
}

func (m *mockUserStore) CreateUser(*types.User) uint {
	return 1
}

func (m *mockUserStore) UpdateUser(*types.User) {}

func (m *mockUserStore) GetRoleByLevel(level int) types.UserRole {
	return types.UserRole{ID: 1}
}

func (m *mockUserStore) AssigneeUserRole(uid uint, role uint) uint {
	return 1
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}
