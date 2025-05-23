package utils

import (
	"userService/models"

	"github.com/go-playground/validator/v10"
)

// for now we are keeping it for user type later it can made generic type
func ValidateStruct(user *models.User) error {
	validate := validator.New()
	return validate.Struct(user)
}
