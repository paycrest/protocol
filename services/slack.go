package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paycrest/aggregator/ent"
	"github.com/paycrest/aggregator/utils"
	"github.com/paycrest/aggregator/utils/logger"
)

type SlackService struct {
	SlackWebhookURL string
}

func NewSlackService(webhookURL string) *SlackService {
	return &SlackService{
		SlackWebhookURL: webhookURL,
	}
}

func (s *SlackService) SendUserSignupNotification(user *ent.User, scopes []string, providerCurrency string) error {
	if s.SlackWebhookURL == "" {
		return fmt.Errorf("slack webhook URL not configured")
	}

	// Format the timestamp using the utility function
	formattedTime, err := utils.FormatTimestampToGMT1(user.CreatedAt)
	if err != nil {
		return fmt.Errorf("error formatting timestamp: %v", err)
	}

	// Prepare Slack message
	message := map[string]interface{}{
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": "*New User Signup*",
				},
			},
			{
				"type": "section",
				"fields": []map[string]interface{}{
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*User ID:* %s", user.ID),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Email:* %s", user.Email),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Name:* %s %s", user.FirstName, user.LastName),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Scopes:* %v", scopes),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Timestamp:* %s", formattedTime),
					},
				},
			},
		},
	}

	// Add provider details if applicable
	if utils.ContainsString(scopes, "provider") && providerCurrency != "" {
		message["blocks"] = append(message["blocks"].([]map[string]interface{}),
			map[string]interface{}{
				"type": "section",
				"fields": []map[string]interface{}{
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Provider Currency:* %s", providerCurrency),
					},
				},
			},
		)
	}

	// Send notification
	jsonPayload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.SlackWebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Errorf("Failed to send Slack notification: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Log error if notification failed, but don't interrupt registration
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Slack notification failed with status: %d", resp.StatusCode)
	}

	return nil
}
