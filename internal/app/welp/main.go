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
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/zlepper/welp/internal/pkg/flatfile"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/templates"
	"strconv"
	"time"
)

type BindWebArgs struct {
	// If the web server should automatically fetch an https certificate and use that
	// If enabled, will ignore the port argument, and always bind on port 80 and 433
	// Will also configure HSTS
	UseHttps bool
	// The port to bind to
	Port int

	// Where to save the uploaded files
	FolderPath string

	// The filename of the flatfile database
	DatabaseFileName string
}

func BindWeb(args BindWebArgs) {
	e := echo.New()

	var logger models.Logger = e.Logger
	fileStorage, err := getFileStorage(args, logger)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	dataStorage, err := getDataStorage(args, logger)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	e.Logger.SetLevel(log.DEBUG)

	setupMiddleware(args, e)

	t := &templateRenderer{
		templates: templates.Must(templates.GetTemplates()),
	}

	e.Renderer = t

	rootGroup := e.Group("")

	bindFeedbackApi(rootGroup, logger, dataStorage, fileStorage)

	if args.UseHttps {
		hostHttps(e, args)
	} else {
		hostHttp(e, args)
	}
}

func hostHttps(e *echo.Echo, args BindWebArgs) {
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}

func hostHttp(e *echo.Echo, args BindWebArgs) {
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(args.Port)))
}

func getFileStorage(args BindWebArgs, logger models.Logger) (models.FileStorage, error) {
	return flatfile.NewFileStorage(flatfile.FileStorageArgs{
		Logger:     logger,
		FolderPath: args.FolderPath,
	})
}

func getDataStorage(args BindWebArgs, logger models.Logger) (models.FeedbackDataStorage, error) {
	return flatfile.NewFeedbackDataStorage(context.Background(), flatfile.DataStorageArgs{
		Logger:       logger,
		SaveInterval: 1 * time.Second,
		Filename:     args.DatabaseFileName,
	})
}

func setupMiddleware(args BindWebArgs, e *echo.Echo) {
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RemoveTrailingSlash(),
		middleware.CORS(),
		middleware.Gzip(),
	)

	if args.UseHttps {
		e.Use(middleware.HTTPSRedirect())
	}
}
