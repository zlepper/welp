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
	"errors"
	"github.com/labstack/echo"
	"github.com/zlepper/welp/internal/app/welp/internal"
	"github.com/zlepper/welp/internal/pkg/consts"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/webapi"
	"net/http"
	"strings"
)

type bindUserManagementApiArgs struct {
	Logger        models.Logger
	DataStorage   models.AuthorizationDataStorage
	AuthService   models.AuthorizationService
	JwtMiddleware echo.MiddlewareFunc
}

func bindUserManagementApi(e *echo.Group, args bindUserManagementApiArgs) {
	server := &userManagementServer{
		bindUserManagementApiArgs: args,
	}

	userGroup := e.Group("/users", args.JwtMiddleware, internal.RequiresRoleMiddleware(models.AdminRole.Key, args.Logger))

	userGroup.GET("", server.getUserList)
	userGroup.POST("", server.postCreateUser)
	userGroup.DELETE("/:email", server.deleteUser)
	userGroup.PUT("/:email", server.updateUser)
	userGroup.GET("/:email", server.getSingleUser)

	userGroup.GET("/new", server.getCreateNewUser)
}

type userManagementServer struct {
	bindUserManagementApiArgs
	baseApi
}

type userListResponse struct {
	AuthState authState
	Users     []models.User
	models.RoleFinder
}

func (r userListResponse) IsLastAdmin(user models.User) bool {
	if !user.HasRole(models.AdminRole.Key) {
		return false
	}

	// Check if there actually are any other admin users
	count := 0
	for _, u := range r.Users {
		if u.HasRole(models.AdminRole.Key) {
			count++
		}
	}

	return count == 1
}

func (s *userManagementServer) getUserList(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

	users, err := s.DataStorage.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	response := userListResponse{
		AuthState: s.getAuthState(c),
		Users:     users,
	}

	return s.respond(c, http.StatusOK, response, "user-list")
}

type createUserRequest struct {
	Name           string   `json:"name" form:"name" xml:"name" query:"name"`
	Email          string   `json:"email" form:"email" xml:"email" query:"email"`
	Password       string   `json:"password" form:"password" xml:"password" query:"password"`
	RepeatPassword string   `json:"repeatPassword" form:"repeatPassword" xml:"repeatPassword" query:"repeatPassword"`
	Roles          []string `json:"roles" form:"roles" xml:"roles" query:"roles"`
}

func (r *createUserRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name required")
	}

	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email required")
	}

	if !strings.ContainsRune(r.Email, '@') {
		return errors.New("email not email")
	}

	if strings.TrimSpace(r.Password) == "" {
		return errors.New("password required")
	}

	if r.Password != r.RepeatPassword {
		return errors.New("password mismatch")
	}

outSearch:
	for _, role := range r.Roles {
		for _, validRole := range models.ValidRoles {
			if role == validRole {
				continue outSearch
			}
		}
		return errors.New("invalid role")
	}

	return nil
}

func (s *userManagementServer) postCreateUser(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

	var request createUserRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = s.createUser(ctx, request)
	if err != nil {
		return err
	}

	return s.respond(c, http.StatusCreated, consts.Nothing, "")
}

func (s *userManagementServer) createUser(ctx context.Context, request createUserRequest) error {

	err := request.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.AuthService.CreateUser(ctx, request.Name, request.Email, request.Password, request.Roles)
	if err != nil {
		return err
	}

	return nil
}

func (s *userManagementServer) deleteUser(c echo.Context) error {
	return nil
}

func (s *userManagementServer) updateUser(c echo.Context) error {
	return nil
}

type errorList []string

func (e errorList) HasError(name string) bool {
	for _, err := range e {
		if err == name {
			return true
		}
	}
	return false
}

type getCreateNewUserData struct {
	AuthState      authState
	AvailableRoles []models.Role
	Errors         errorList
}

func (s *userManagementServer) getCreateNewUser(c echo.Context) error {
	response := getCreateNewUserData{
		AuthState:      s.getAuthState(c),
		AvailableRoles: models.Roles,
	}
	return s.respond(c, http.StatusOK, response, "create-new-user")
}

func (s *userManagementServer) postCreateNewUser(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

	var request createUserRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = s.createUser(ctx, request)
	if err != nil {
		response := getCreateNewUserData{
			AvailableRoles: models.Roles,
			AuthState:      s.getAuthState(c),
			Errors:         []string{err.Error()},
		}
		return c.Render(http.StatusBadRequest, "create-new-user", response)
	}

	return c.Redirect(http.StatusSeeOther, "/users")
}

type singleUserResponse struct {
	AuthState authState
	User      models.User
	models.RoleFinder
	AvailableRoles []models.Role
}

func (s *userManagementServer) getSingleUser(c echo.Context) error {
	ctx := webapi.GetContext(c.Request())

	email := c.Param("email")

	user, err := s.DataStorage.GetUser(ctx, email)
	if err != nil {
		if err == models.ErrNoSuchUser {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return err
	}

	res := singleUserResponse{
		User:           user,
		AuthState:      s.getAuthState(c),
		AvailableRoles: models.Roles,
	}

	return s.respond(c, http.StatusOK, res, "edit-user")
}
