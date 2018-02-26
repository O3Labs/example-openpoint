package services

import (
	"fmt"

	"github.com/o3labs/openpoint/platform/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/o3labs/openpoint/platform/config"
)

type EmailService struct {
	EmailServiceInterface
}

type EmailServiceInterface interface {
	SendEmail(email models.Email) error
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (e *EmailService) SendEmail(email models.Email) error {

	sess := session.New(config.AWSConfig())

	svc := ses.New(sess, nil)
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: []*string{},
			CcAddresses:  []*string{},
			ToAddresses: []*string{
				aws.String(email.To),
			},
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
				Html: &ses.Content{
					Data: aws.String(email.HTML), // Required
				},
			},
			Subject: &ses.Content{
				Data: aws.String(email.Subject), // Required
			},
		},
		Source: aws.String(email.From), // Required
		ReplyToAddresses: []*string{
			aws.String(email.From), // Required
		},
	}
	resp, err := svc.SendEmail(params)

	if err != nil {
		return err
	}
	fmt.Printf(resp.String())
	return nil
}
