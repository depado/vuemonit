package models

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const bcost = 12

type User struct {
	ID       string `storm:"id" json:"id"`
	Email    string `storm:"unique" json:"email"`
	Password []byte `json:"password"`

	CreatedAt time.Time `storm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}

func NewUser(email, password string) (*User, error) {
	var err error
	var pwd []byte

	if pwd, err = bcrypt.GenerateFromPassword([]byte(password), bcost); err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	return &User{
		ID:        xid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Email:     email,
		Password:  pwd,
	}, nil
}

// CheckPassword is a simple utility function to check the password given as raw
// against the user's hashed password
func (u User) CheckPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(raw)) == nil
}
