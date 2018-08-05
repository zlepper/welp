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

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/zlepper/welp/internal/pkg/consts"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/webapi"
	"net/http"
	"strings"
)

var (
	ErrAuthorizationRequired = errors.New("authorization required")
)

type UnauthorizedResponse struct {
	Message string `xml:"message" json:"message"`
}

func GetJWTMiddlware(secretService models.SecretService, logger models.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			signingSecret, err := secretService.GetSigningSecret(webapi.GetContext(c.Request()))
			if err != nil {
				return err
			}

			var user models.TokenUser
			err = getTokenDataFromRequest(c.Request(), signingSecret, &user)
			if err != nil {

				logger.Infof("Authorization required: %v", err)

				responseType := webapi.GetResponseType(c.Request())

				returnUrl := c.Request().URL.EscapedPath()

				logger.Debugf("Returning to: %s", returnUrl)

				switch responseType {
				case webapi.MIMEJSON:
				case webapi.MIMEXML:
				default:
					return echo.ErrUnauthorized
				case webapi.MIMEHTML:
					return c.Redirect(http.StatusSeeOther, "/login?returnUrl="+returnUrl)
				}
			}

			c.Set("user", user)

			return next(c)
		}
	}
}

func getValidationKeyGetter(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", method)
			}
		}
		return secret, nil
	}
}

func getTokenData(tokenString string, secret []byte, output interface{}) error {
	token, err := jwt.Parse(tokenString, getValidationKeyGetter(secret))
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub := claims["sub"].(string)
		err = json.Unmarshal([]byte(sub), output)
		return err
	} else {
		return err
	}
}

func getTokenDataFromRequest(request *http.Request, secret []byte, output interface{}) error {
	var tokenString string

	// Try from header
	authHeader := request.Header.Get(echo.HeaderAuthorization)
	if authHeader != "" {
		if !strings.HasPrefix(authHeader, consts.Bearer) {
			return ErrAuthorizationRequired
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return ErrAuthorizationRequired
		}

		tokenString = parts[1]
	} else {
		// Try query parameters
		tokenString = request.URL.Query().Get("token")
		if tokenString == "" {
			// Try from cookies
			cookie, err := request.Cookie(echo.HeaderAuthorization)
			if err != nil {
				return ErrAuthorizationRequired
			}
			tokenString = cookie.Value
		}

		if tokenString == "" {
			return ErrAuthorizationRequired
		}
	}

	return getTokenData(tokenString, secret, output)

}
