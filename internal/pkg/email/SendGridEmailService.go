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

package email

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/zlepper/welp/internal/pkg/models"
)

type SendGridEmailServiceArgs struct {
	logger models.Logger
	// The sendgrid api key
	ApiKey string
}

func NewSendGridEmailService(args SendGridEmailServiceArgs) models.EmailService {
	return &sendGridEmailService{
		logger: args.logger,
		apiKey: args.ApiKey,
	}
}

type sendGridEmailService struct {
	logger models.Logger
	apiKey string
}

func (s *sendGridEmailService) SendEmail(args models.SendEmailArgs) error {

	from := mail.NewEmail(args.From.Name, args.From.Address)
	var replyTo *mail.Email
	if args.ReplyTo.Address == "" {
		replyTo = from
	} else {
		replyTo = mail.NewEmail(args.ReplyTo.Name, args.ReplyTo.Address)
	}

	subject := args.Subject
	to := mail.NewEmail(args.To.Name, args.To.Address)
	plainTextContent := args.PlainContent
	htmlContent := args.HtmlContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	message.SetReplyTo(replyTo)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil

}
