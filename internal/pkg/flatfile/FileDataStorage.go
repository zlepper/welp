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
	"encoding/json"
	"github.com/zlepper/welp/internal/pkg/models"
	"os"
	"sync"
	"time"
)

type DataStorageArgs struct {
	// The name of the file to save data to
	Filename string

	// How often the data changes should be saved
	SaveInterval time.Duration

	// A logger for logging information
	Logger models.Logger
}

// A simple file storage. Just saves things to a local flatfile
// Keeps everything in memory meanwhile
func NewFeedbackDataStorage(ctx context.Context, args DataStorageArgs) (models.FeedbackDataStorage, error) {
	storage := &feedbackFileDataStorage{
		args:    args,
		changed: false,
		data:    map[string]models.Feedback{},
		logger:  args.Logger,
	}

	err := storage.load()
	if err != nil {
		args.Logger.Errorf("Failed to load stored data: %v", err)
		return nil, err
	}

	go storage.startSaveCycle(ctx)

	return storage, nil
}

func (s *feedbackFileDataStorage) startSaveCycle(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Shutting down file data storage")
			s.saveChanges()
			// We are completely done saving
			return
		case <-time.After(s.args.SaveInterval):
			err := s.saveChanges()
			if err != nil {
				s.logger.Error(err)
			}
			// And just start waiting again
		}
	}

}

func (s *feedbackFileDataStorage) saveChanges() error {
	// If there are no changes to save at all, then don't do anything
	if !s.changed {
		return nil
	}

	s.logger.Info("Data has changed. Saving...")

	s.lock.Lock()
	defer s.lock.Unlock()

	tempName := s.args.Filename + ".temp"

	// Save to temp file first to ensure no errors
	file, err := os.Create(tempName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(s.data)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tempName, s.args.Filename)
	if err != nil {
		return err
	}

	// Register that there are no new changes to be saved
	s.changed = false

	return nil
}

// Loads the current data from disk
func (s *feedbackFileDataStorage) load() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	file, err := os.Open(s.args.Filename)
	if err != nil {
		// If the file doesn't exist, then it's simply because we haven't saved anything yet.
		if os.IsNotExist(err) {
			s.logger.Infof("Data file '%s' doesn't exist. Assuming it's because it hasn't been created yet.", s.args.Filename)
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&s.data)
}

type feedbackFileDataStorage struct {
	// The args for the data storage
	args DataStorageArgs

	// Indicates if anything has changed, since the last time the data was saved
	changed bool

	// The actual data
	data map[string]models.Feedback

	// Lock to ensure exclusivity
	lock sync.RWMutex

	logger models.Logger
}

func (s *feedbackFileDataStorage) SaveFeedback(ctx context.Context, feedback models.Feedback) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.changed = true

	s.data[feedback.Id] = feedback

	return nil
}
