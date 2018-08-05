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
	"golang.org/x/crypto/acme/autocert"
	"strconv"
)

func host(args models.BindWebArgs, e *echo.Echo) {

	if args.UseHttps {
		hostHttps(e, args)
	} else {
		hostHttp(e, args)
	}
}

func hostHttps(e *echo.Echo, args models.BindWebArgs) {
	e.AutoTLSManager.Cache = autocert.DirCache(args.CertificateCacheFolder)
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}

func hostHttp(e *echo.Echo, args models.BindWebArgs) {
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(args.Port)))
}
