package notifiers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"website-monitor/result"
)

type SlackNotifier struct {
	name       string
	webhookUrl string
}

var SlackMissingWebhookErr = errors.New("required option 'webhook' is missing")

func NewSlackNotifier(name string, options map[string]string) (*SlackNotifier, error) {
	if _, ok := options["webhook"]; !ok {
		return nil, SlackMissingWebhookErr
	}

	sn := &SlackNotifier{
		name:       name,
		webhookUrl: options["webhook"],
	}

	delete(options, "webhook")

	return sn, nil
}

type SlackTextSection struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackBlock struct {
	Type string           `json:"type"`
	Text SlackTextSection `json:"text"`
}

type SlackRequestBody struct {
	Text   string       `json:"text"`
	Blocks []SlackBlock `json:"blocks"`
}

func (s *SlackNotifier) Notify(name, displayUrl string, result *result.Results) error {
	var text string
	if result.AllTrue() {
		text = fmt.Sprintf("<%s|%s> *matches* checks!", displayUrl, name)
	} else {
		text = fmt.Sprintf("%s does *not* match checks!", name)
	}

	body := SlackRequestBody{}
	body.Text = text
	body.Blocks = append(body.Blocks, SlackBlock{
		Type: "section",
		Text: SlackTextSection{
			Type: "mrkdwn",
			Text: text,
		},
	})
	for _, r := range result.Results {
		body.Blocks = append(body.Blocks, SlackBlock{
			Type: "section",
			Text: SlackTextSection{
				Type: "mrkdwn",
				Text: fmt.Sprintf("%s: %t (err: %v)", r.ContentChecker, r.Result, r.Err),
			},
		})
	}

	slackBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, s.webhookUrl, bytes.NewBuffer(slackBody))
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

func (s *SlackNotifier) Name() string {
	return s.name
}

func (s *SlackNotifier) Equal(y *SlackNotifier) bool {
	if s.name != y.name {
		return false
	}

	return s.webhookUrl == y.webhookUrl
}
