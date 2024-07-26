package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	mailgunv3 "github.com/mailgun/mailgun-go/v3"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/types"
)

var (
	notificationConf = config.NotificationConfig()

	mailgunClient       mailgunv3.Mailgun
	_DefaultFromAddress = notificationConf.EmailFromAddress //"Paycrest <no-reply@paycrest.io>"
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

// NewMailgun initialize mailgunv3.Mailgun and can be used to initialize a mocked Mailgun interface.
func NewMailgun(m mailgunv3.Mailgun) {
	if m != nil {
		mailgunClient = m
		return
	}

	mailgunClient = mailgunv3.NewMailgun(notificationConf.EmailDomain, notificationConf.EmailAPIKey)
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
	payload := types.SendEmailPayload{
		FromAddress: _DefaultFromAddress,
		ToAddress:   email,
		DynamicData: map[string]interface{}{
			"code": token,
		},
	}
	return SendTemplateEmail(payload, "d-f26d853bbb884c0c856f0bbda894032c")

}

// SendPasswordResetEmail performs the actions for sending a password reset token to the user email.
func (m *EmailService) SendPasswordResetEmail(ctx context.Context, token, email string) (types.SendEmailResponse, error) {

	payload := types.SendEmailPayload{
		FromAddress: _DefaultFromAddress,
		ToAddress:   email,
		DynamicData: map[string]interface{}{
			"code": token,
		},
	}
	return SendTemplateEmail(payload, "d-8b689801cd9947748775ccd1c4cc932e")
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

func SendTemplateEmail(content types.SendEmailPayload, templateId string) (types.SendEmailResponse, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := "https://api.sendgrid.com/v3/mail/send"

	reqBody := map[string]interface{}{
		"from": map[string]string{
			"email": _DefaultFromAddress,
		},
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{
						"email": content.ToAddress,
					},
				},
				"dynamic_template_data": content.DynamicData,
			},
		},
		"template_id": templateId,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return types.SendEmailResponse{}, fmt.Errorf("error marshalling JSON: %w", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return types.SendEmailResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+notificationConf.EmailAPIKey)

	body, err := client.Do(req)
	if err != nil {
		return types.SendEmailResponse{}, fmt.Errorf("error sending request: %w", err)
	}
	return types.SendEmailResponse{
		Response: body.Header.Get("X-Message-Id"),
		Id:       body.Header.Get("X-Message-Id"),
	}, nil
}

func SendTemplateEmailWithJsonAttachment(content types.SendEmailPayload, templateId string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := "https://api.sendgrid.com/v3/mail/send"

	reqBody := map[string]interface{}{
		"from": map[string]string{
			"email": _DefaultFromAddress,
		},
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{
						"email": content.ToAddress,
					},
				},
				"dynamic_template_data": content.DynamicData,
			},
		},
		"template_id": templateId,
		"attachments": []map[string]interface{}{
			{
				"content": content.Body,
				"type":    "text/json", "disposition": "attachment",
				"filename": "payload.json",
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+notificationConf.EmailAPIKey)

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	return nil
}
