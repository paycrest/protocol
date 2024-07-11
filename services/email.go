package services

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/types"
)

var (
	notificationConf = config.NotificationConfig()

	mailgunClient mailgun.Mailgun
)

const (
	_DefaultFromAddress        = "Paycrest <no-reply@paycrest.io>"
	_EmailVerificationSubject  = "Your Paycrest Email Verification Token"
	_PasswordResetEmailSubject = "Reset Password confirmation"
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

// NewMailgun initialize mailgun.Mailgun and can be used to initialize a mocked Mailgun interface.
func NewMailgun(m mailgun.Mailgun) {
	if m != nil {
		mailgunClient = m
		return
	}

	mailgunClient = mailgun.NewMailgun(notificationConf.EmailDomain, notificationConf.EmailAPIKey)
}

// SendEmail performs the action for sending an email.
func (m *EmailService) SendEmail(ctx context.Context, payload types.SendEmailPayload) (types.SendEmailResponse, error) {
	switch m.MailProvider {
	case MAILGUN_MAIL_PROVIDER:
		return sendEmailViaMailgun(ctx, payload)
	case SENDGRID_MAIL_PROVIDER:
		return sendEmailViaSendGrid(ctx, payload)
	default:
		return types.SendEmailResponse{}, fmt.Errorf("unsupported mail provider")
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
		HTMLBody:    fmt.Sprintf("<div> %s </div>", bodyTemplate),
	}
	return m.SendEmail(ctx, payload)
}

// SendPasswordResetEmail performs the actions for sending a password reset token to the user email.
func (m *EmailService) SendPasswordResetEmail(ctx context.Context, token, email string) (types.SendEmailResponse, error) {
	// TODO: add custom HTML email verification template.
	body := fmt.Sprintf("Please confirm your password reset request with this token: %s", token)
	htmlBody := ""

	payload := types.SendEmailPayload{
		FromAddress: _DefaultFromAddress,
		ToAddress:   email,
		Subject:     _PasswordResetEmailSubject,
		Body:        body,
		HTMLBody:    htmlBody,
	}
	return m.SendEmail(ctx, payload)
}

// sendEmailViaMailgun performs the actions for sending an email.
func sendEmailViaMailgun(ctx context.Context, content types.SendEmailPayload) (types.SendEmailResponse, error) {
	// initialize
	NewMailgun(mailgunClient)

	message := mailgunClient.NewMessage(
		content.FromAddress,
		content.Subject,
		content.Body,
		content.ToAddress,
	)

	response, id, err := mailgunClient.Send(ctx, message)

	return types.SendEmailResponse{Id: id, Response: response}, err
}

// sendEmailViaSendGrid performs the actions for sending an email.
func sendEmailViaSendGrid(ctx context.Context, content types.SendEmailPayload) (types.SendEmailResponse, error) {
	_ = ctx
	from := mail.NewEmail("Paycrest", "<no-reply@paycrest.io>")
	to := mail.NewEmail("", content.ToAddress)
	body := mail.NewContent("text/plain", content.Body)
	htmlBody := mail.NewContent("text/html", content.HTMLBody)

	m := mail.NewV3Mail()
	m.Subject = content.Subject
	m.SetFrom(from)
	m.AddContent(body)
	m.AddContent(htmlBody)

	p := mail.NewPersonalization()
	p.AddTos(to)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(notificationConf.EmailAPIKey, "/v3/mail/send", fmt.Sprintf("https://%s", notificationConf.EmailDomain))
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil || response.StatusCode >= 400 {
		return types.SendEmailResponse{}, err
	}

	return types.SendEmailResponse{Id: response.Headers["X-Message-Id"][0]}, nil
}

// sendEmailViaSendGrid performs the actions for sending an email.
func sendEmailViaTemplates(ctx context.Context, content types.SendEmailPayload) (types.SendEmailResponse, error) {
	// initialize
	NewMailgun(mailgunClient)

	// The message object allows you to add attachments and Bcc recipients
	message := mailgunClient.NewMessage("<no-reply@paycrest.io>", content.Subject, "", content.ToAddress)
	message.SetTemplate("No yet given")
	err := message.AddTemplateVariable("c1", "value")
	if err != nil {
		return types.SendEmailResponse{}, err
	}
	err = message.AddTemplateVariable("c2", "value2")
	if err != nil {
		return types.SendEmailResponse{}, err
	}

	err = message.AddTemplateVariable("c3", "value")
	if err != nil {
		return types.SendEmailResponse{}, err
	}

	err = message.AddTemplateVariable("c4", "value")
	if err != nil {
		return types.SendEmailResponse{}, err
	}

	err = message.AddTemplateVariable("c5", "value")
	if err != nil {
		return types.SendEmailResponse{}, err
	}

	err = message.AddTemplateVariable("c6", "value")
	if err != nil {
		return types.SendEmailResponse{}, err
	}


	// Send the message with a 10 second timeout
	response, id, err := mailgunClient.Send(ctx, message)

	if err != nil {
		return types.SendEmailResponse{}, err
	}

	return types.SendEmailResponse{Id: id, Response: response}, err
}

func GetTemplateIds(ctx context.Context) *mailgun.TemplatesIterator {
	mailgunClient := mailgun.NewMailgun(notificationConf.EmailDomain, notificationConf.EmailAPIKey)
	return mailgunClient.ListTemplates(&mailgun.ListTemplateOptions{Limit: 20, Active: true})
}
