package services

import (
	"context"
	"fmt"

	mailgunv3 "github.com/mailgun/mailgun-go/v3"
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/types"
)

var (
	emailConf = config.EmailConfig()

	mailGunClient mailgunv3.Mailgun
)

const (
	_DefaultFromAddress       = "Paycrest <no-reply@paycrest.io>"
	_EmailVerificationSubject = "Your Paycrest Email Verification Token"
)

type MailProvider string

const (
	MAILGUN_MAIL_PROVIDER  MailProvider = "MAILGUN"
	SENDGRID_MAIL_PROVIDER MailProvider = "SENDGRID"
)

// EmailService provides functionality to sending emails via a mailer provider
type EmailService struct {
	MailProvider MailProvider
}

// NewEmailService creates a new instance of EmailService with a given MailProvider.
func NewEmailService(mailProvider MailProvider) *EmailService {
	return &EmailService{MailProvider: mailProvider}
}

// SendEmail performs the action for sending a email
func (m *EmailService) SendEmail(ctx context.Context, payload types.SendEmailPayload) (types.SendEmailResponse, error) {
	switch m.MailProvider {
	case MAILGUN_MAIL_PROVIDER:
		fallthrough
	default:
		return sendEmailViaMailgun(ctx, payload)
	}
}

// SendVerificationEmail performs the actions for sending a verification token to the user email.
func (m *EmailService) SendVerificationEmail(ctx context.Context, token, email string) (types.SendEmailResponse, error) {
	// TODO: add custom HTML email verification template.
	bodyTemplate := fmt.Sprintf("Verification Token: %s", token)

	payload := types.SendEmailPayload{
		FromAddress: _DefaultFromAddress,
		ToAddress:   email,
		Subject:     _EmailVerificationSubject,
		Body:        bodyTemplate,
	}

	return m.SendEmail(ctx, payload)
}

// NewMailGun initialize mailgunv3.Mailgun and can be used to initialize a mocked Mailgun interface.
func NewMailGun(m mailgunv3.Mailgun) {
	if m != nil {
		mailGunClient = m
		return
	}

	mailGunClient = mailgunv3.NewMailgun(emailConf.Domain, emailConf.ApiKey)
}

// sendEmailViaMailgun performs the actions for sending an email.
func sendEmailViaMailgun(ctx context.Context, content types.SendEmailPayload) (types.SendEmailResponse, error) {
	// initialize
	NewMailGun(mailGunClient)

	message, id, err := mailGunClient.Send(ctx, mailGunClient.NewMessage(
		content.FromAddress,
		content.Subject,
		content.Body,
		content.ToAddress,
	))

	return types.SendEmailResponse{Id: id, Message: message}, err
}
