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
	"github.com/zlepper/welp/internal/pkg/consts"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/webapi"
	"net/http"
	"time"
)

type AuthorizationApiArgs struct {
	Logger        models.Logger
	AuthService   models.AuthorizationService
	LoginDuration time.Duration
	JwtMiddleware echo.MiddlewareFunc
}

type authorizationApiServer struct {
	baseApi
	AuthorizationApiArgs
}

func bindAuthorizationApi(g *echo.Group, args AuthorizationApiArgs) {
	authApiServer := &authorizationApiServer{
		AuthorizationApiArgs: args,
	}

	g.GET("/login", authApiServer.loginGetHandler)
	g.POST("/login", authApiServer.loginPostHandler)
	g.GET("/logout", authApiServer.logoutGetHandler)
}

type getLoginResponse struct {
	AuthState authState
}

func (s *authorizationApiServer) loginGetHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login", getLoginResponse{
		AuthState: s.getAuthState(c),
	})
}

type loginRequest struct {
	Email    string `json:"email" form:"email" xml:"email" query:"email"`
	Password string `json:"password" form:"password" xml:"password" query:"password"`
}

type loginResponse struct {
	Token string `json:"token" xml:"token"`
}

func (s *authorizationApiServer) loginPostHandler(c echo.Context) error {
	var request loginRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	token, err := s.AuthService.Login(webapi.GetContext(c.Request()), request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	responseType := webapi.GetResponseType(c.Request())

	c.SetCookie(&http.Cookie{
		Name:    echo.HeaderAuthorization,
		Domain:  c.Request().URL.Host,
		Expires: time.Now().Add(s.LoginDuration),
		Path:    "/",
		Value:   token,
	})

	switch responseType {
	case webapi.MIMEHTML:
	default:
		returnUrl := c.Param("returnUrl")

		if returnUrl == "" {
			returnUrl = "/"
		}

		return c.Redirect(http.StatusSeeOther, returnUrl)
	case webapi.MIMEJSON:
		return c.JSON(http.StatusOK, loginResponse{Token: token})
	case webapi.MIMEXML:
		return c.JSON(http.StatusOK, loginResponse{Token: token})
	}

	return nil
}

func (s *authorizationApiServer) logoutGetHandler(c echo.Context) error {

	// Send an expired cookie, to remove the current logged in cookie
	c.SetCookie(&http.Cookie{
		Name:    echo.HeaderAuthorization,
		Domain:  c.Request().URL.Host,
		Expires: time.Unix(0, 0),
		Path:    "/",
		Value:   "",
	})

	responseType := webapi.GetResponseType(c.Request())
	switch responseType {
	case webapi.MIMEJSON:
		return c.JSON(http.StatusOK, consts.Nothing)
	case webapi.MIMEXML:
		return c.XML(http.StatusOK, consts.Nothing)
	default:
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}
