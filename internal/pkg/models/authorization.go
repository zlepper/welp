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
)

type Role struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Roles    []Role `json:"roles"`
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
}
