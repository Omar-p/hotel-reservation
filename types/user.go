package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 8
)

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u UpdateUserRequest) Validate() map[string]string {
	errors := map[string]string{}
	if u.FirstName != "" && len(u.FirstName) < minFirstNameLength {
		errors["firstName"] = "firstName must be at least 2 characters"
	}
	if u.LastName != "" && len(u.LastName) < minLastNameLength {
		errors["lastName"] = "lastName must be at least 2 characters"
	}
	return errors
}

func (u UpdateUserRequest) ToBSON() bson.M {
	m := bson.M{}
	if u.FirstName != "" {

		m["firstName"] = u.FirstName
	}
	if u.LastName != "" {
		m["lastName"] = u.LastName
	}
	return m
}
func (c CreateUserRequest) Validate() map[string]string {
	errors := map[string]string{}
	if len(c.FirstName) < minFirstNameLength {
		errors["firstName"] = "firstName must be at least 2 characters"
	}
	if len(c.LastName) < minLastNameLength {
		errors["lastName"] = "lastName must be at least 2 characters"
	}
	if len(c.Password) < minPasswordLength {
		errors["password"] = "password must be at least 8 characters"
	}
	if !isValidEmail(c.Email) {
		errors["email"] = "email must be a valid email"
	}
	return errors
}

func isValidEmail(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"lastName"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromCreateRequest(req *CreateUserRequest) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		EncryptedPassword: string(encryptedPassword),
	}, nil
}
