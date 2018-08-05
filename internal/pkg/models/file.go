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
	"io"
	"strings"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

// How to save and load actual attached files
type FileStorage interface {
	// Should save the file with the given name
	// The file content can be read from the reader
	SaveFile(ctx context.Context, id string, reader io.Reader) (size int64, err error)

	// Should load the given file name from disk/storage
	LoadFile(ctx context.Context, id string) (reader io.ReadCloser, err error)
}

type File struct {
	// The id used to refer to the file
	Id string `json:"id"`
	// The total size of the file
	Size int64 `json:"size"`
	// The contentType of the file
	ContentType string `json:"contentType"`
}

// Returns true if this file is actually an image
func (f *File) IsImage() bool {
	return strings.HasPrefix(f.ContentType, "image/")
}
