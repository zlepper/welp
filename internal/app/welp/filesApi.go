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
	"github.com/labstack/echo"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/webapi"
	"io"
	"math"
	"net/http"
	"strconv"
)

type filesApiArgs struct {
	models.Logger
	models.FeedbackDataStorage
	models.FileStorage
	JwtMiddleware echo.MiddlewareFunc
}

func bindFilesApi(e *echo.Group, args filesApiArgs) {

	server := &filesApiServer{
		filesApiArgs: args,
		baseApi:      baseApi{},
	}

	e.GET("/files/:id", server.getFileHandler, args.JwtMiddleware)
}

type filesApiServer struct {
	filesApiArgs
	baseApi
}

type getFileNotFound struct {
	Message string `json:"message" xml:"message"`
}

func (s *filesApiServer) getFileHandler(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

	id := c.Param("id")

	reader, err := s.FileStorage.LoadFile(ctx, id)
	if err != nil {
		if err == models.ErrFileNotFound {
			return s.respond(c, http.StatusNotFound, getFileNotFound{Message: err.Error()}, "file-not-found")
		}
		return err
	}
	defer reader.Close()

	c.Response().Header().Set(webapi.HeaderCacheControl, "public, max-age="+strconv.Itoa(math.MaxInt32))

	_, err = io.Copy(c.Response(), reader)
	return err
}
