package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*
	Package helpers provides utility functions for password hashing, user authorization, and password verification.

	This package includes the following functions:

	- HashPassword: Hashes the provided password using bcrypt. It returns the hashed password as a string and an error if any.

	- CheckUserType: Checks if the user type in the context matches the provided user role. It returns an error if the user is not authorized to access the resource.

	- VerifyPassword: Compares the hashed password with the input password. It returns true if the passwords match, false otherwise.
*/

// HashPassword hashes the provided password using bcrypt.
// It returns the hashed password as a string and an error if any.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckUserType checks if the user type in the context matches the provided user role.
// It returns an error if the user is not authorized to access the resource.
func CheckUserType(c *gin.Context, userRole string) error {
	userType := c.GetString("user_type")
	if userType != userRole {
		return fmt.Errorf("user is not authorized to access this resource")
	}
	return nil
}

// VerifyPassword compares the hashed password with the input password.
// It returns true if the passwords match, false otherwise.
func VerifyPassword(hashedPassword string, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
