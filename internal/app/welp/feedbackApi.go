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
	"github.com/zlepper/welp/internal/pkg/webapi"
	"mime/multipart"
	"net/http"
	"path"
)

type feedbackServer struct {
	bindFeedbackApiArgs
	baseApi
}

type bindFeedbackApiArgs struct {
	Logger          models.Logger
	FeedbackService models.FeedbackService
	FileStorage     models.FileStorage
	EmailService    models.EmailService
	JwtMiddleware   echo.MiddlewareFunc
}

func bindFeedbackApi(e *echo.Group, args bindFeedbackApiArgs) {

	server := &feedbackServer{
		bindFeedbackApiArgs: args,
	}

	e.POST("/", server.createFeedbackEntryHandler)
	e.GET("/embed", server.getFeedbackEmbedHandler)
	e.GET("/", server.getFeedbackListHandler, args.JwtMiddleware)
}

type createFeedbackRequest struct {
	// The actual message from the user
	Message string `json:"message" form:"message"`
	// A contact email address, in case the user want to be informed once their feedback
	// has been addressed
	ContactAddress string `json:"contactAddress" form:"contactAddress"`
}

func (s *feedbackServer) createFeedbackEntryHandler(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

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

	savedFiles := make([]models.File, 0)

	for _, file := range files {
		if file.Size == 0 {
			continue
		}

		created, err := s.saveMultipartFile(ctx, file)
		if err != nil {
			return err
		}

		savedFiles = append(savedFiles, created)
	}

	feedback, err := s.FeedbackService.CreateFeedback(ctx, request.Message, request.ContactAddress, savedFiles)
	if err != nil {
		return err
	}

	return s.respond(c, http.StatusCreated, feedback, "feedback-created")
}

func (s *feedbackServer) sendNewFeedbackEmail(ctx context.Context, feedback models.Feedback) {

}

func (s *feedbackServer) saveMultipartFile(ctx context.Context, file *multipart.FileHeader) (createdFile models.File, err error) {

	src, err := file.Open()
	if err != nil {
		return createdFile, err
	}
	defer src.Close()
	contentType, reader, err := webapi.DetectContentType(src)

	id, err := uuid.NewRandom()
	if err != nil {
		return createdFile, err
	}

	ext := path.Ext(file.Filename)

	filename := id.String() + ext

	size, err := s.FileStorage.SaveFile(ctx, filename, reader)
	if err != nil {
		return createdFile, err
	}

	return models.File{
		Id:          filename,
		Size:        size,
		ContentType: contentType,
	}, nil
}

func (s *feedbackServer) getFeedbackEmbedHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "embed", nil)
}

func (s *feedbackServer) getAllFeedbackHandler(c echo.Context) error {
	feedback, err := s.FeedbackService.GetAllFeedback(context.Background())
	if err != nil {
		return err
	}

	return s.respond(c, http.StatusOK, feedback, "feedback-list")
}

type feedbackResponse struct {
	Feedback  []models.Feedback
	AuthState authState
}

func (s *feedbackServer) getFeedbackListHandler(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

	feedback, err := s.FeedbackService.GetAllFeedback(ctx)
	if err != nil {
		return err
	}

	response := feedbackResponse{
		Feedback:  feedback,
		AuthState: s.getAuthState(c),
	}

	return s.respond(c, http.StatusOK, response, "feedback-list")
}
