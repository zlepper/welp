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
	"io"
	"os"
	"path"
)

type FileStorageArgs struct {
	// The path of the folder to put all files into
	FolderPath string

	Logger models.Logger
}

// Gets a new file storage for saving files to the local disk
func NewFileStorage(args FileStorageArgs) (models.FileStorage, error) {
	storage := &fileStorage{
		args:   args,
		logger: args.Logger,
	}

	err := os.MkdirAll(args.FolderPath, os.ModePerm)
	if err != nil {
		args.Logger.Errorf("Failed to create file storage location: %v", err)
		return nil, err
	}

	return storage, nil
}

type fileStorage struct {
	args   FileStorageArgs
	logger models.Logger
}

func (s *fileStorage) getPath(name string) string {
	return path.Join(s.args.FolderPath, name)
}

func (s *fileStorage) SaveFile(ctx context.Context, name string, reader io.Reader) (size int64, err error) {
	filename := s.getPath(name)
	s.logger.Debugf("Saving file to: '%s'", filename)

	file, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, reader)
}

func (s *fileStorage) LoadFile(ctx context.Context, name string) (reader io.ReadCloser, err error) {
	filename := s.getPath(name)
	s.logger.Debugf("Loading file from '%s'", filename)

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, models.ErrFileNotFound
		}
		return nil, err
	}

	return file, nil
}
