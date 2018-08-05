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

package flatfile

import (
	"context"
	"github.com/zlepper/welp/internal/pkg/models"
	"sync"
	"time"
)

type AuthorizationDataStorageArgs struct {
	// The name of the file to save data to
	Filename string

	// How often the data changes should be saved
	SaveInterval time.Duration

	// A logger for logging information
	Logger models.Logger
}

func NewAuthorizationDataStorage(ctx context.Context, args AuthorizationDataStorageArgs) (models.AuthorizationDataStorage, error) {
	storage := &authorizationDataStorage{
		data:    map[string]models.User{},
		changed: false,
		logger:  args.Logger,
	}

	saver := NewDataSaver(DataSaverArgs{
		logger:       args.Logger,
		saveInterval: args.SaveInterval,
		saveable:     storage,
		filename:     args.Filename,
	})

	err := saver.LoadData(&storage.data)
	if err != nil {
		args.Logger.Errorf("Failed to load stored data: %v", err)
		return nil, err
	}

	go saver.StartSaveCycle(ctx)

	return storage, nil
}

type authorizationDataStorage struct {
	lock    sync.RWMutex
	changed bool
	logger  models.Logger
	data    map[string]models.User
}

func (s *authorizationDataStorage) GetUser(ctx context.Context, email string) (models.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	user, exists := s.data[email]
	if exists {
		return user, nil
	}

	return models.User{}, models.ErrNoSuchUser
}

func (s *authorizationDataStorage) CreateUser(ctx context.Context, user models.User) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.data[user.Email]; exists {
		return models.ErrUserAlreadyExists
	}

	s.data[user.Email] = user

	s.changed = true

	return nil
}

func (s *authorizationDataStorage) GetAllUsers(ctx context.Context) ([]models.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	out := make([]models.User, len(s.data))
	index := 0
	for _, user := range s.data {
		out[index] = user
	}

	return out, nil
}

func (s *authorizationDataStorage) DeleteUser(ctx context.Context, email string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.data, email)

	s.changed = true

	return nil
}

func (s *authorizationDataStorage) GetUserCount(ctx context.Context) (int, error) {
	return len(s.data), nil
}

func (s *authorizationDataStorage) Lock() {
	s.lock.Lock()
}

func (s *authorizationDataStorage) Unlock() {
	s.lock.Unlock()
}

func (s *authorizationDataStorage) GetData() interface{} {
	return s.data
}

func (s *authorizationDataStorage) HasChanged() bool {
	return s.changed
}

func (s *authorizationDataStorage) SetChanged(changed bool) {
	s.changed = changed
}
