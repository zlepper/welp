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
)

type EmailSender struct {
	FromName, FromEmail       string
	ReplyToName, ReplyToEmail string
}

type AuthorizationServiceArgs struct {
	Logger       models.Logger
	DataStorage  models.AuthorizationDataStorage
	EmailService models.EmailService
	EmailSender  EmailSender
}

func NewAuthorizationService(ctx context.Context, args AuthorizationServiceArgs) (*AuthorizationService, error) {
	return &AuthorizationService{
		logger:       args.Logger,
		dataStorage:  args.DataStorage,
		emailService: args.EmailService,
		emailSender:  args.EmailSender,
	}, nil
}

type AuthorizationService struct {
	logger       models.Logger
	dataStorage  models.AuthorizationDataStorage
	emailService models.EmailService
	emailSender  EmailSender
}

func (s *AuthorizationService) hashPassword(password string) (hash string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(hash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(h), nil
}

func (s *AuthorizationService) comparePasswords(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Creates a new user
// roles should be a slice of the keys of the roles the user should have
func (s *AuthorizationService) CreateUser(ctx context.Context, email string, password string, roles []string) error {

	hash, err := s.hashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: hash,
		Roles:    []models.Role{},
	}

	err = s.dataStorage.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	err = s.SendWelcomeEmail(email)
	return err
}

func (s *AuthorizationService) SendWelcomeEmail(emailAddress string) error {
	return s.emailService.SendEmail(models.SendEmailArgs{
		From:         models.NewEmailAddress(s.emailSender.FromName, s.emailSender.FromEmail),
		ReplyTo:      models.NewEmailAddress(s.emailSender.ReplyToName, s.emailSender.ReplyToEmail),
		To:           models.NewEmailAddress(emailAddress, emailAddress),
		Subject:      "Welcome to Welp",
		PlainContent: "Welcome to Welp",
		HtmlContent:  "Welcome to Welp",
	})
}
