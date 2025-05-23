package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
	"userService/models"

	"github.com/google/uuid"
)

var userMap = make(map[uuid.UUID]*models.User)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, update *models.UpdateUserModel) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserServiceImpl struct {
	Users map[uuid.UUID]*models.User
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{userMap}
}

// CreateUser will create a user in database
func (u *UserServiceImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	log.Printf("service started with Create User")
	defer log.Printf("service ends with Create User")

	if user == nil {
		return nil, errors.New("user model is empty")
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	u.Users[user.ID] = user
	return user, nil

}

// GetUser will fetch a user from database
func (u *UserServiceImpl) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	log.Printf("service started with Get User")
	defer log.Printf("service ends with Get User")

	if len(strings.TrimSpace(id.String())) == 0 {
		return nil, errors.New("empty id provided")
	}

	// check if user is present in database
	user, present := u.Users[id]
	if present {
		return user, nil
	}

	return nil, errors.New("user not present")

}

// UpdateUser will update a user from database
func (u *UserServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, update *models.UpdateUserModel) (*models.User, error) {
	log.Printf("service started with update User")
	defer log.Printf("service ends with update User")

	if len(strings.TrimSpace(id.String())) == 0 {
		return nil, errors.New("empty id provided")
	}

	// check if user is present in database
	user, present := u.Users[id]
	if !present {
		return nil, errors.New("user not present")
	}

	if len(update.Name) != 0 {
		user.Name = update.Name
	}

	if len(update.Address) != 0 {
		user.Address = update.Address
	}

	user.UpdatedAt = time.Now()
	u.Users[id] = user

	return user, nil

}

// DeleteUser will delete a user from database
func (u *UserServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	log.Printf("service started with delete User")
	defer log.Printf("service ends with delete User")

	if len(strings.TrimSpace(id.String())) == 0 {
		return errors.New("empty id provided")
	}

	// check if user is present in database
	_, present := u.Users[id]
	if !present {
		return errors.New("user not present")
	}

	delete(u.Users, id)

	log.Println("user successfully deleted")

	return nil

}
