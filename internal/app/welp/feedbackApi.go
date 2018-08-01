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

package welp

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/zlepper/welp/internal/pkg/models"
	"mime/multipart"
	"net/http"
	"path"
)

type feedbackServer struct {
	logger      models.Logger
	dataStorage models.FeedbackDataStorage
	fileStorage models.FileStorage
	baseApi
}

func bindFeedbackApi(e *echo.Group, logger models.Logger, dataStorage models.FeedbackDataStorage, fileStorage models.FileStorage) {

	g := e.Group("/feedback")

	server := &feedbackServer{
		logger:      logger,
		dataStorage: dataStorage,
		fileStorage: fileStorage,
	}

	g.POST("", server.createFeedbackEntryHandler)
	g.GET("/embed", server.getFeedbackEmbedHandler)
}

type createFeedbackRequest struct {
	// The actual message from the user
	Message string `json:"message" form:"message"`
	// A contact email address, in case the user want to be informed once their feedback
	// has been addressed
	ContactAddress string `json:"contactAddress" form:"contactAddress"`
}

func (s *feedbackServer) createFeedbackEntryHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var request createFeedbackRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// Read all the files and save them to storage
	files := form.File["files"]

	savedFiles := make([]models.File, len(files))

	for index, file := range files {
		created, err := s.saveMultipartFile(ctx, file)
		if err != nil {
			return err
		}

		savedFiles[index] = created
	}

	feedback, err := models.NewFeedback(request.Message, request.ContactAddress, savedFiles)
	if err != nil {
		return err
	}

	s.logger.Infof("Saving changes")
	err = s.dataStorage.SaveFeedback(ctx, feedback)
	if err != nil {
		return err
	}

	return s.respond(c, http.StatusCreated, feedback, "feedback-created")
}

func (s *feedbackServer) saveMultipartFile(ctx context.Context, file *multipart.FileHeader) (createdFile models.File, err error) {
	src, err := file.Open()
	if err != nil {
		return createdFile, err
	}
	defer src.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		return createdFile, err
	}

	ext := path.Ext(file.Filename)

	filename := id.String() + ext

	size, err := s.fileStorage.SaveFile(ctx, filename, src)
	if err != nil {
		return createdFile, err
	}

	return models.File{
		Id:       id.String(),
		Location: filename,
		Name:     file.Filename,
		Size:     size,
	}, nil
}

func (s *feedbackServer) getFeedbackEmbedHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "embed", nil)
}
