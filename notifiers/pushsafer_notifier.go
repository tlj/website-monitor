package notifiers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"website-monitor/result"
)

type PushSaferNotifier struct {
	privateKey string
	options    map[string]string
}

var PushSaferMissingPrivateKeyErr = errors.New("required option 'private_key' is missing")

func NewPushSaferNotifier(options map[string]string) (*PushSaferNotifier, error) {
	if _, ok := options["private_key"]; !ok {
		return nil, PushSaferMissingPrivateKeyErr
	}

	delete(options, "private_key")

	return &PushSaferNotifier{
		privateKey: options["private_key"],
		options:    options,
	}, nil
}

func (p PushSaferNotifier) Notify(name, displayUrl string, result *result.Results) error {
	params := url.Values{}

	for k, v := range p.options {
		params.Set(k, v)
	}

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
	if p.privateKey != y.privateKey {
		return false
	}

	return reflect.DeepEqual(p.options, y.options)
}
