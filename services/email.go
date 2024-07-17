package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	mailgunv3 "github.com/mailgun/mailgun-go/v3"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
)

var (
	notificationConf = config.NotificationConfig()

	mailgunClient mailgunv3.Mailgun
)

const (
	_DefaultFromAddress        = "oayobami15@gmail.com" //"Paycrest <no-reply@paycrest.io>"
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

func SendTemplateEmail(content types.SendEmailPayload, templateId string) error {
	client := fastshot.NewClient("https://api.sendgrid.com/").
		Config().SetTimeout(30*time.Second).
		Header().Set("Content-Type", "application/json").
		Auth().BearerToken(notificationConf.EmailAPIKey)

	res, err := client.Build().POST("/v3/mail/send").
		Body().AsJSON(map[string]interface{}{
		"from": map[string]string{
			"email": content.FromAddress,
		},
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{
						"email": content.ToAddress,
					},
				},
				"dynamic_template_data": content.DynamicData,
				"subject":               content.Subject,
			},
		},
		"subject":     content.Subject,
		"template_id": templateId,
	}).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		return fmt.Errorf("fetch txn event logs: %v, %v", err, res.RawResponse)
	}

	_, err = utils.ParseJSONResponse(res.RawResponse)
	if err != nil {
		return fmt.Errorf("fetch: %v, %v", err, res.RawResponse)
	}

	return nil
}

func SendTemplateEmailWithJsonAttachment(content types.SendEmailPayload, templateId string) error {
	// client := fastshot.NewClient("https://api.sendgrid.com/v3/mail/send").
	// 	Config().SetTimeout(30*time.Second).
	// 	Header().Set("Content-Type", "application/json").
	// 	Auth().BearerToken(notificationConf.EmailAPIKey)

	// ok := utils.IsBase64(content.Body)
	// if !ok {
	// 	return fmt.Errorf("attachments is not Base64")
	// }
	// res, err := client.Build().POST("").
	// 	Body().AsJSON(map[string]interface{}{
	// 	"from": map[string]string{
	// 		"email": _DefaultFromAddress,
	// 	},
	// 	"personalizations": []map[string]interface{}{
	// 		{
	// 			"to": []map[string]string{
	// 				{
	// 					"email": content.ToAddress,
	// 				},
	// 			},
	// 			"dynamic_template_data": content.DynamicData,
	// 		},
	// 	},
	// 	"template_id": templateId,
	// 	// "attachments": []map[string]interface{}{
	// 	// 	{
	// 	// 		"content":     content.Body,
	// 	// 		"type":        "text/json",
	// 	// 		"disposition": "attachment",
	// 	// 		"filename":    "payload.json",
	// 	// 	},
	// 	// },
	// }).
	// 	// Retry().Set(3, 1*time.Second).
	// 	Send()
	// if err != nil {
	// 	return fmt.Errorf("sendEmail: %v,v", err, res.RawResponse)
	// }

	// _, err = utils.ParseJSONResponse(res.RawResponse)
	// if err != nil {

	// 	return fmt.Errorf("failed to Parse: %v, %v", err, res.RawResponse)
	// }

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
		return fmt.Errorf("Error marshalling JSON: %w", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+notificationConf.EmailAPIKey)

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %w", err)
	}

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return fmt.Errorf("Error reading response body:", err)
	// }

	// return fmt.Errorf("Error reading response body: %v,\n %v \n %v", string(body), content.DynamicData, templateId)

	// fmt.Println("Response status:", resp.Status)
	// fmt.Println("Response body:", string(body))
	return nil
}
func SendVerificationEmailV2(content types.SendEmailPayload) error {
	// TODO: Not sure how you want the templateId added yet
	return SendTemplateEmail(content, "d-be6f6e592fd242ca9e14db589f21c1b1")
}
