package services

import (
	"context"
	"github.com/zlepper/welp/internal/pkg/models"
)

type FeedbackServiceArgs struct {
	DataStorage     models.FeedbackDataStorage
	EmailService    models.EmailService
	UserDataStorage models.AuthorizationDataStorage
	Logger          models.Logger
	Args            models.BindWebArgs
}

func NewFeedbackService(args FeedbackServiceArgs) models.FeedbackService {
	return &feedbackService{
		FeedbackServiceArgs: args,
	}
}

type feedbackService struct {
	FeedbackServiceArgs
}

func (s *feedbackService) CreateFeedback(ctx context.Context, message, contactAddress string, files []models.File) (models.Feedback, error) {
	feedback, err := models.NewFeedback(message, contactAddress, files)
	if err != nil {
		return models.Feedback{}, err
	}

	err = s.DataStorage.SaveFeedback(ctx, feedback)
	if err != nil {
		return models.Feedback{}, err
	}

	go s.sendFeedbackEmails(ctx, feedback)

	return feedback, nil
}

func (s *feedbackService) GetAllFeedback(ctx context.Context) ([]models.Feedback, error) {
	return s.DataStorage.GetAllFeedback(ctx)
}

func (s *feedbackService) sendFeedbackEmails(ctx context.Context, feedback models.Feedback) error {
	users, err := s.UserDataStorage.GetAllUsers(ctx)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	var from models.EmailAddress
	if feedback.ContactAddress == "" {
		from = models.NewEmailAddress(s.Args.EmailSenderName, s.Args.EmailSenderAddress)
	} else {
		from = models.NewEmailAddress("User", feedback.ContactAddress)
	}

	for _, user := range users {
		if user.EmailUpdate == models.Immediately {
			err = s.EmailService.SendEmail(models.SendEmailArgs{
				Subject:      "New feedback",
				From:         from,
				ReplyTo:      from,
				To:           models.NewEmailAddress(user.Name, user.Email),
				PlainContent: feedback.Message,
				HtmlContent:  feedback.Message,
			})
			if err != nil {
				s.Logger.Error(err)
				return err
			}
		}
	}

	s.Logger.Info("Emails send to all people who wanted immediate feedback")

	return nil
}
