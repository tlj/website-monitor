package notifiers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type SlackNotifier struct {
	WebhookUrl string
}

func NewSlackNotifier(webhookUrl string) *SlackNotifier {
	return &SlackNotifier{
		WebhookUrl: webhookUrl,
	}
}

func (s *SlackNotifier) Notify(msg string) error {
	return sendSlackNotification(s.WebhookUrl, msg)
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func sendSlackNotification(webhookUrl string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if buf.String() != "ok" {
		return errors.New("non-ok response returned from Slack")
	}

	return nil
}