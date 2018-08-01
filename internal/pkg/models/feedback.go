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
	"github.com/google/uuid"
)

type FeedbackDataStorage interface {
	SaveFeedback(ctx context.Context, feedback Feedback) error
}

func NewFeedback(message, contactAddress string, files []File) (Feedback, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Feedback{}, err
	}

	return Feedback{
		Id:             id.String(),
		ContactAddress: contactAddress,
		Message:        message,
		Files:          files,
	}, nil
}

type Feedback struct {
	// The id of the feedback entry
	Id string `json:"id"`
	// The message attached to the feedback
	Message string `json:"message"`
	// Files that was attached to the feedback
	Files []File `json:"files"`
	// A contact address for getting back to the user who provided the feedback
	ContactAddress string `json:"contactAddress"`
}
