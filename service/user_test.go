package service

import (
	"context"
	"testing"
	"time"
	"userService/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	ctx := context.Background()
	svc := NewUserService()

	t.Run("CreateUser - Success", func(t *testing.T) {
		user := &models.User{Name: "Alice", Address: "123 Wonderland"}
		createdUser, err := svc.CreateUser(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.NotEqual(t, uuid.Nil, createdUser.ID)
		assert.WithinDuration(t, time.Now(), createdUser.CreatedAt, time.Second)
	})

	t.Run("CreateUser - Nil User", func(t *testing.T) {
		createdUser, err := svc.CreateUser(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, createdUser)
	})

	t.Run("GetUser - Success", func(t *testing.T) {
		user := &models.User{Name: "Bob", Address: "456 Builder Blvd"}
		createdUser, _ := svc.CreateUser(ctx, user)

		fetchedUser, err := svc.GetUser(ctx, createdUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, createdUser, fetchedUser)
	})

	t.Run("GetUser - Not Found", func(t *testing.T) {
		_, err := svc.GetUser(ctx, uuid.New())

		assert.Error(t, err)
	})

	t.Run("UpdateUser - Success", func(t *testing.T) {
		user := &models.User{Name: "Charlie", Address: "789 Chocolate Factory"}
		createdUser, _ := svc.CreateUser(ctx, user)

		update := &models.UpdateUserModel{Name: "Charlie Updated", Address: "New Address"}
		updatedUser, err := svc.UpdateUser(ctx, createdUser.ID, update)

		assert.NoError(t, err)
		assert.Equal(t, "Charlie Updated", updatedUser.Name)
		assert.Equal(t, "New Address", updatedUser.Address)
		assert.WithinDuration(t, time.Now(), updatedUser.UpdatedAt, time.Second)
	})

	t.Run("UpdateUser - Not Found", func(t *testing.T) {
		update := &models.UpdateUserModel{Name: "No User"}
		_, err := svc.UpdateUser(ctx, uuid.New(), update)

		assert.Error(t, err)
	})

	t.Run("DeleteUser - Success", func(t *testing.T) {
		user := &models.User{Name: "Delete Me", Address: "Gone Soon"}
		createdUser, _ := svc.CreateUser(ctx, user)

		err := svc.DeleteUser(ctx, createdUser.ID)
		assert.NoError(t, err)

		// Confirm user is gone
		_, err = svc.GetUser(ctx, createdUser.ID)
		assert.Error(t, err)
	})

	t.Run("DeleteUser - Not Found", func(t *testing.T) {
		err := svc.DeleteUser(ctx, uuid.New())
		assert.Error(t, err)
	})
}
