package mailgun

import (
	"context"
	"fmt"
	"time"

	mailgunv3 "github.com/mailgun/mailgun-go/v3"
	"github.com/paycrest/paycrest-protocol/config"
)

var mailgunConf = config.MailGunConfig()

type MailGunResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

func sendMessage(from, subject, text, to string) (MailGunResponse, error) {
	mg := mailgunv3.NewMailgun(mailgunConf.Domain, mailgunConf.ApiKey)
	message := mg.NewMessage(from, subject, text, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	m, id, err := mg.Send(ctx, message)

	return MailGunResponse{ID: id, Message: m}, err
}

func SendVerificationEmail(token, to string) (MailGunResponse, error) {
	return sendMessage("Paycrest <no-reply@paycrest.io>", "Your Paycrest Email Verification Token", fmt.Sprintf("Verification Token: %s", token), to)
}
