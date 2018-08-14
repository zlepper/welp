/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package models

import (
	"context"
	"errors"
	"log"
)

var (
	AdminRole       = Role{"Admin", "admin"}
	ErrRoleNotFound = errors.New("role not found")
)

var Roles = []Role{AdminRole}

func FindRole(key string) (Role, error) {
	for _, role := range Roles {
		if role.Key == key {
			return role, nil
		}
	}
	return Role{}, ErrRoleNotFound
}

type RoleFinder struct {
}

func (r RoleFinder) GetRole(key string) Role {
	role, err := FindRole(key)
	if err != nil {
		log.Panicln(err)
	}

	return role
}

var ValidRoles = []string{
	AdminRole.Key,
}

type Role struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type EmailNotificationUpdate string

const (
	// The user should never receive any emails
	Never EmailNotificationUpdate = "never"
	// The user should receive a daily update with all the feedback that came that day,
	// in one big email
	Daily EmailNotificationUpdate = "daily"
	// The user should receiver an email with the new feedback the moment it comes in
	Immediately EmailNotificationUpdate = "immediately"
)

type User struct {
	Name        string                  `json:"name"`
	Email       string                  `json:"email"`
	Password    string                  `json:"password"`
	Roles       []string                `json:"roles"`
	EmailUpdate EmailNotificationUpdate `json:"emailUpdate"`
}

func (u User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

// A user variant, that can be passed around in tokens
type TokenUser struct {
	Email string
	// The role keys of the roles the user has
	Roles []string
}

func (u TokenUser) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoSuchUser        = errors.New("no such user")
)

type AuthorizationDataStorage interface {
	// Should fetch the user with the given username from the datastorage
	// If the user doesn't exist, ErrNoSuchUser should be returned
	GetUser(ctx context.Context, email string) (User, error)
	// Should create the specified user if possible
	// If the user already exists, then method should error
	CreateUser(ctx context.Context, user User) error
	// Should get all the users currently in the data storage
	GetAllUsers(ctx context.Context) ([]User, error)
	// Should delete the user with the given username
	DeleteUser(ctx context.Context, email string) error
	// Should get the number of users in the system
	GetUserCount(ctx context.Context) (int, error)
	// Updates an existing user in the system
	UpdateUser(ctx context.Context, email string, user User) error
}

type AuthorizationService interface {
	CreateUser(ctx context.Context, name, email, password string, roles []string) error
	Login(ctx context.Context, email, password string) (string, error)
}
