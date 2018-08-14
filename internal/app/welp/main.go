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
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/zlepper/welp/internal/app/welp/internal"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/templates"
)

func BindWeb(args models.BindWebArgs) {
	e := echo.New()

	e.HideBanner = true

	var logger models.Logger = e.Logger

	loadedServices, err := internal.GetServices(args, logger)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	e.Logger.SetLevel(log.DEBUG)

	setupMiddleware(args, e)
	jwtMiddleware := internal.GetJWTMiddlware(loadedServices.SecretService, logger)

	t := &templateRenderer{
		templates: templates.Must(templates.GetTemplates()),
	}

	e.Renderer = t

	rootGroup := e.Group("")

	bindFeedbackApi(rootGroup, bindFeedbackApiArgs{
		FileStorage:     loadedServices.FileStorage,
		Logger:          logger,
		FeedbackService: loadedServices.FeedbackService,
		JwtMiddleware:   jwtMiddleware,
	})

	bindAuthorizationApi(rootGroup, AuthorizationApiArgs{
		Logger:        logger,
		AuthService:   loadedServices.AuthorizationService,
		LoginDuration: args.TokenDuration,
		JwtMiddleware: jwtMiddleware,
	})

	bindFilesApi(rootGroup, filesApiArgs{
		JwtMiddleware:       jwtMiddleware,
		FileStorage:         loadedServices.FileStorage,
		FeedbackDataStorage: loadedServices.FeedbackDataStorage,
		Logger:              logger,
	})

	bindUserManagementApi(rootGroup, bindUserManagementApiArgs{
		Logger:        logger,
		JwtMiddleware: jwtMiddleware,
		DataStorage:   loadedServices.AuthorizationDataStorage,
		AuthService:   loadedServices.AuthorizationService,
	})

	host(args, e)
}

func setupMiddleware(args models.BindWebArgs, e *echo.Echo) {
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RemoveTrailingSlash(),
		middleware.CORS(),
		middleware.Gzip(),
	)

	if args.UseHttps {
		e.Pre(middleware.HTTPSRedirect())
	}
}
