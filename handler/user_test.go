package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"userService/models"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UpdateUserModel) (*models.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUserHandler(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandlerService(mockSvc)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{Name: "John Doe", Address: "address example asdf"}
		userJson, _ := json.Marshal(user)

		mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/users/v1", bytes.NewBuffer(userJson))
		rec := httptest.NewRecorder()

		h.CreateUserHandler(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockSvc.AssertCalled(t, "CreateUser", mock.Anything, mock.AnythingOfType("*models.User"))
	})

	t.Run("Invalid Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users/v1", bytes.NewBuffer([]byte(`{invalid json`)))
		rec := httptest.NewRecorder()

		h.CreateUserHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		user := &models.User{Name: "John Doe", Address: "address example asdf"}
		userJson, _ := json.Marshal(user)

		mockSvc.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil, errors.New("internal error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/users/v1", bytes.NewBuffer(userJson))
		rec := httptest.NewRecorder()

		h.CreateUserHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateUserHandler(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandlerService(mockSvc)

	validUUID := uuid.New()
	updateModel := &models.UpdateUserModel{Name: "Updated Name"}
	updatedUser := &models.User{Name: "Updated Name"}
	updateJson, _ := json.Marshal(updateModel)

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("UpdateUser", mock.Anything, validUUID, updateModel).Return(updatedUser, nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/"+validUUID.String(), bytes.NewBuffer(updateJson))
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.UpdateUserHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockSvc.AssertCalled(t, "UpdateUser", mock.Anything, validUUID, updateModel)
	})

	t.Run("Missing ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/", bytes.NewBuffer(updateJson))
		rec := httptest.NewRecorder()

		h.UpdateUserHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/invalid-uuid", bytes.NewBuffer(updateJson))
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
		rec := httptest.NewRecorder()

		h.UpdateUserHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Decode Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/"+validUUID.String(), bytes.NewBuffer([]byte("{bad json")))
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.UpdateUserHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.On("UpdateUser", mock.Anything, validUUID, updateModel).Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/"+validUUID.String(), bytes.NewBuffer(updateJson))
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.UpdateUserHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetUserHandler(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandlerService(mockSvc)

	validUUID := uuid.New()
	user := &models.User{Name: "John Doe"}

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("GetUser", mock.Anything, validUUID).Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/"+validUUID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.GetUserHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Missing ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		rec := httptest.NewRecorder()

		h.GetUserHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
		rec := httptest.NewRecorder()

		h.GetUserHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.On("GetUser", mock.Anything, validUUID).Return(nil, errors.New("not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/"+validUUID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.GetUserHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDeleteUserHandler(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandlerService(mockSvc)

	validUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockSvc.On("DeleteUser", mock.Anything, validUUID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/"+validUUID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.DeleteUser(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Missing ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/", nil)
		rec := httptest.NewRecorder()

		h.DeleteUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
		rec := httptest.NewRecorder()

		h.DeleteUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.On("DeleteUser", mock.Anything, validUUID).Return(errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/"+validUUID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
		rec := httptest.NewRecorder()

		h.DeleteUser(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
