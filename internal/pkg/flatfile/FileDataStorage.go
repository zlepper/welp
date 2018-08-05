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
	"sort"
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

	saver := NewDataSaver(DataSaverArgs{
		filename:     args.Filename,
		saveable:     storage,
		logger:       args.Logger,
		saveInterval: args.SaveInterval,
	})

	err := saver.LoadData(&storage.data)
	if err != nil {
		args.Logger.Errorf("Failed to load stored data: %v", err)
		return nil, err
	}

	go saver.StartSaveCycle(ctx)

	return storage, nil
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

func (s *feedbackFileDataStorage) Lock() {
	s.lock.Lock()
}

func (s *feedbackFileDataStorage) Unlock() {
	s.lock.Unlock()
}

func (s *feedbackFileDataStorage) GetData() interface{} {
	return s.data
}

func (s *feedbackFileDataStorage) HasChanged() bool {
	return s.changed
}

func (s *feedbackFileDataStorage) SetChanged(changed bool) {
	s.changed = changed
}

func (s *feedbackFileDataStorage) SaveFeedback(ctx context.Context, feedback models.Feedback) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.changed = true

	s.data[feedback.Id] = feedback

	return nil
}

func (s *feedbackFileDataStorage) GetAllFeedback(ctx context.Context) ([]models.Feedback, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	out := make([]models.Feedback, len(s.data))
	index := 0
	for _, value := range s.data {
		out[index] = value
		index++
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Created.After(out[j].Created)
	})

	return out, nil
}
