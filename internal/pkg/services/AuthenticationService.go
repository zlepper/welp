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

package services

import (
	"context"
	"github.com/zlepper/welp/internal/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type EmailSender struct {
	FromName, FromEmail       string
	ReplyToName, ReplyToEmail string
}

//noinspection GoNameStartsWithPackageName
type AuthorizationServiceArgs struct {
	Logger        models.Logger
	DataStorage   models.AuthorizationDataStorage
	EmailService  models.EmailService
	EmailSender   EmailSender
	TokenService  models.TokenService
	TokenDuration time.Duration
}

func NewAuthorizationService(args AuthorizationServiceArgs) (models.AuthorizationService, error) {
	return &authorizationService{
		logger:        args.Logger,
		dataStorage:   args.DataStorage,
		emailService:  args.EmailService,
		emailSender:   args.EmailSender,
		tokenService:  args.TokenService,
		tokenDuration: args.TokenDuration,
	}, nil
}

type authorizationService struct {
	logger        models.Logger
	dataStorage   models.AuthorizationDataStorage
	emailService  models.EmailService
	emailSender   EmailSender
	tokenService  models.TokenService
	tokenDuration time.Duration
}

func (s *authorizationService) hashPassword(password string) (hash string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(h), nil
}

func (s *authorizationService) comparePasswords(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Creates a new user
// roles should be a slice of the keys of the roles the user should have
func (s *authorizationService) CreateUser(ctx context.Context, name, email, password string, roles []string) error {

	// Verify roles exists
outerSearch:
	for _, r := range roles {
		for _, role := range models.Roles {
			if role.Key == r {
				continue outerSearch
			}
		}
		return models.ErrRoleNotFound
	}

	hash, err := s.hashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:        name,
		Email:       email,
		Password:    hash,
		Roles:       roles,
		EmailUpdate: models.Never,
	}

	err = s.dataStorage.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	err = s.sendWelcomeEmail(email)
	return err
}

func (s *authorizationService) sendWelcomeEmail(emailAddress string) error {
	// TODO Better welcome email
	return s.emailService.SendEmail(models.SendEmailArgs{
		From:         models.NewEmailAddress(s.emailSender.FromName, s.emailSender.FromEmail),
		ReplyTo:      models.NewEmailAddress(s.emailSender.ReplyToName, s.emailSender.ReplyToEmail),
		To:           models.NewEmailAddress(emailAddress, emailAddress),
		Subject:      "Welcome to Welp",
		PlainContent: "Welcome to Welp",
		HtmlContent:  "Welcome to Welp",
	})
}

func (s *authorizationService) ensureAtLeastOneUserExists(ctx context.Context) error {
	count, err := s.dataStorage.GetUserCount(ctx)
	if err != nil {
		return err
	}

	if count != 0 {
		// User already exists, so we don't need to do any more
		return nil
	}

	hash, err := s.hashPassword("admin")
	if err != nil {
		return err
	}

	user := models.User{
		Email:    "admin@admin.com",
		Password: hash,
		Roles:    []string{models.AdminRole.Key},
	}

	err = s.dataStorage.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authorizationService) Login(ctx context.Context, email, password string) (string, error) {
	err := s.ensureAtLeastOneUserExists(ctx)
	if err != nil {
		return "", nil
	}

	user, err := s.dataStorage.GetUser(ctx, email)
	if err != nil {
		return "", err
	}

	s.logger.Debugf("pw: %s, hash: %s", password, user.Password)
	err = s.comparePasswords(password, user.Password)
	if err != nil {
		return "", err
	}

	token, err := s.tokenService.GenerateToken(ctx, s.tokenDuration, s.generateTokenUser(user))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authorizationService) generateTokenUser(user models.User) models.TokenUser {
	return models.TokenUser{
		Email: user.Email,
		Roles: user.Roles,
	}
}
