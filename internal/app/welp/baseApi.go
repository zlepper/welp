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
	"net/http"
)

type baseApi struct {
}

// Responds to the request in the way to client prefers
func (b *baseApi) respond(c echo.Context, code int, response interface{}, templateName string) error {
	responseType := webapi.GetResponseType(c.Request())

	switch responseType {
	case webapi.MIMEJSON:
	default:
		return c.JSON(http.StatusCreated, response)
	case webapi.MIMEXML:
		return c.XML(http.StatusCreated, response)
	case webapi.MIMEHTML:
		return c.Render(http.StatusCreated, templateName, response)
	}

	return nil
}

type authState struct {
	Authenticated bool
	User          models.TokenUser
}

func (b *baseApi) getAuthState(c echo.Context) authState {
	user := c.Get("user")
	if user != nil {
		tokenUser := user.(models.TokenUser)
		return authState{
			Authenticated: true,
			User:          tokenUser,
		}
	}

	return authState{
		Authenticated: false,
	}
}
