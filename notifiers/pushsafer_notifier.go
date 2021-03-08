package notifiers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"website-monitor/result"
)

type PushSaferNotifier struct {
	privateKey string
}

func NewPushSaferNotifier(privateKey string) *PushSaferNotifier {
	return &PushSaferNotifier{
		privateKey: privateKey,
	}
}

func (p PushSaferNotifier) Notify(name, displayUrl string, result *result.Results) error {
	params := url.Values{}

	params.Set("k", p.privateKey)
	params.Set("t", name)
	params.Set("u", displayUrl)
	params.Set("d", "a") // all devices

	if result.AllTrue() {
		params.Set("m", fmt.Sprintf("<%s|%s> *matches* checks!", displayUrl, name))
	} else {
		params.Set("m", fmt.Sprintf("%s does *not* match checks!", name))
	}

	res, err := http.Post("https://www.pushsafer.com/api", "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code from pushsafer: %d", res.StatusCode)
	}

	return nil
}

func (p *PushSaferNotifier) Equal(y *PushSaferNotifier) bool {
	return p.privateKey == y.privateKey
}
