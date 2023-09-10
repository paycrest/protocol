package mailgun

import (
	"context"
	"fmt"
	"time"

	mailgunv3 "github.com/mailgun/mailgun-go/v3"
	"github.com/paycrest/paycrest-protocol/config"
)

var (
	mailgunConf = config.MailGunConfig()

	mailer mailgunv3.Mailgun
)

// NewMailGun initialize mailgunv3.Mailgun and can be used to initialize a mocked Mailgun interface.
func NewMailGun(m mailgunv3.Mailgun) {
	if _, ok := m.(mailgunv3.Mailgun); ok {
		mailer = m
		return
	}

	mailer = mailgunv3.NewMailgun(mailgunConf.Domain, mailgunConf.ApiKey)
}

// MailGunResponse is the mailgunv3.Send response struct
type MailGunResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

// SendVerificationEmail performs the actions for sending a verification token to the user email.
func SendVerificationEmail(token, email string) (MailGunResponse, error) {
	// initialize
	NewMailGun(mailer)

	from := "Paycrest <no-reply@paycrest.io>"
	subject := "Your Paycrest Email Verification Token"
	// TODO: add custom HTML email verification template.
	text := fmt.Sprintf("Verification Token: %s", token)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	message, id, err := mailer.Send(ctx, mailer.NewMessage(from, subject, text, email))

	return MailGunResponse{Id: id, Message: message}, err
}
