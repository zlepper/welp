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
	"context"
	"github.com/zlepper/welp/internal/pkg/email"
	"github.com/zlepper/welp/internal/pkg/flatfile"
	"github.com/zlepper/welp/internal/pkg/models"
	"github.com/zlepper/welp/internal/pkg/services"
	"path"
)

type dataLayerType int

const (
	unknown           dataLayerType = iota
	flatFileDataLayer dataLayerType = iota
)

type loadedServices struct {
	models.FileStorage
	models.FeedbackDataStorage
	models.SecretService
	models.AuthorizationService
	models.AuthorizationDataStorage
	models.EmailService
	models.FeedbackService
}

func GetServices(args models.BindWebArgs, logger models.Logger) (*loadedServices, error) {

	fileStorage, err := getFileStorage(args, logger)
	if err != nil {
		return nil, err
	}

	feedbackDataStorage, err := getDataStorage(args, logger)
	if err != nil {
		return nil, err
	}

	secretService, err := getSecretService(args, logger)
	if err != nil {
		return nil, err
	}

	emailService, err := getEmailService(args, logger)
	if err != nil {
		return nil, err
	}

	tokenService, err := getTokenService(args, logger, secretService)
	if err != nil {
		return nil, err
	}

	authenticationDataStorage, err := getAuthenticationDataStorage(args, logger)
	if err != nil {
		return nil, err
	}

	authenticationService, err := getAuthenticationService(args, logger, emailService, tokenService, authenticationDataStorage)
	if err != nil {
		return nil, err
	}

	feedbackService, err := getFeedbackService(args, logger, emailService, feedbackDataStorage, authenticationDataStorage)
	if err != nil {
		return nil, err
	}

	return &loadedServices{
		FileStorage:              fileStorage,
		FeedbackDataStorage:      feedbackDataStorage,
		SecretService:            secretService,
		AuthorizationService:     authenticationService,
		AuthorizationDataStorage: authenticationDataStorage,
		EmailService:             emailService,
		FeedbackService:          feedbackService,
	}, nil

}

func detectDataLayerType(args models.BindWebArgs) dataLayerType {
	if args.DatabaseFolderName != "" {
		return flatFileDataLayer
	}
	return unknown
}

func getFileStorage(args models.BindWebArgs, logger models.Logger) (models.FileStorage, error) {
	return flatfile.NewFileStorage(flatfile.FileStorageArgs{
		Logger:     logger,
		FolderPath: args.FolderPath,
	})
}

func getDataStorage(args models.BindWebArgs, logger models.Logger) (models.FeedbackDataStorage, error) {
	return flatfile.NewFeedbackDataStorage(context.Background(), flatfile.DataStorageArgs{
		Logger:       logger,
		SaveInterval: args.SaveInterval,
		Filename:     path.Join(args.DatabaseFolderName, "feedback.json"),
	})
}

func getSecretService(args models.BindWebArgs, logger models.Logger) (models.SecretService, error) {
	return flatfile.NewSecretStorage(flatfile.SecretStorageArgs{
		Logger:   logger,
		Filename: path.Join(args.DatabaseFolderName, "secrets.json"),
	})
}

func getEmailService(args models.BindWebArgs, logger models.Logger) (models.EmailService, error) {
	if args.SendGridApiKey != "" {
		return email.NewSendGridEmailService(email.SendGridEmailServiceArgs{
			ApiKey: args.SendGridApiKey,
			Logger: logger,
		})
	}

	return email.NewNoOpEmailService(), nil
}

func getTokenService(args models.BindWebArgs, logger models.Logger, secretService models.SecretService) (models.TokenService, error) {
	return services.NewTokenService(services.TokenServiceArgs{
		SecretService: secretService,
		Logger:        logger,
	})
}

func getAuthenticationDataStorage(args models.BindWebArgs, logger models.Logger) (models.AuthorizationDataStorage, error) {
	return flatfile.NewAuthorizationDataStorage(context.Background(), flatfile.AuthorizationDataStorageArgs{
		Logger:       logger,
		Filename:     path.Join(args.DatabaseFolderName, "authentication.json"),
		SaveInterval: args.SaveInterval,
	})
}

func getAuthenticationService(args models.BindWebArgs, logger models.Logger, emailService models.EmailService, tokenService models.TokenService, dataStorage models.AuthorizationDataStorage) (models.AuthorizationService, error) {
	return services.NewAuthorizationService(services.AuthorizationServiceArgs{
		Logger:        logger,
		EmailService:  emailService,
		TokenDuration: args.TokenDuration,
		TokenService:  tokenService,
		EmailSender: services.EmailSender{
			FromEmail:    args.EmailSenderAddress,
			FromName:     args.EmailSenderName,
			ReplyToEmail: args.EmailSenderAddress,
			ReplyToName:  args.EmailSenderName,
		},
		DataStorage: dataStorage,
	})
}

func getFeedbackService(args models.BindWebArgs, logger models.Logger, emailService models.EmailService, feedbackDataStorage models.FeedbackDataStorage, userDataStorage models.AuthorizationDataStorage) (models.FeedbackService, error) {
	return services.NewFeedbackService(services.FeedbackServiceArgs{
		Logger:          logger,
		EmailService:    emailService,
		DataStorage:     feedbackDataStorage,
		UserDataStorage: userDataStorage,
		Args:            args,
	}), nil
}
