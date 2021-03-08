package notifiers

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"website-monitor/result"
)

type PushSaferNotifier struct {
	name       string
	privateKey string
	options    map[string]string
}

var PushSaferMissingPrivateKeyErr = errors.New("required option 'private_key' is missing")

func NewPushSaferNotifier(name string, options map[string]string) (*PushSaferNotifier, error) {
	if _, ok := options["private_key"]; !ok {
		return nil, PushSaferMissingPrivateKeyErr
	}

	ps := &PushSaferNotifier{
		name:       name,
		privateKey: options["private_key"],
		options:    options,
	}

	delete(options, "private_key")

	return ps, nil
}

type PushSaferResponse struct {
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Success   string `json:"success"`
	Available int    `json:"available"`
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

	body, _ := ioutil.ReadAll(res.Body)
	jr := PushSaferResponse{}
	err = json.Unmarshal(body, &jr)
	if err != nil {
		return fmt.Errorf("error decoding response from pushsafer: %s", string(body))
	}

	if jr.Status == 0 {
		return fmt.Errorf("error returned from pushsafer: %s", jr.Error)
	}

	log.Debugf("PushSafer '%s' available calls: %d", p.name, jr.Available)

	return nil
}

func (p *PushSaferNotifier) Name() string {
	return p.name
}

func (p *PushSaferNotifier) Equal(y *PushSaferNotifier) bool {
	if p.name != y.name {
		return false
	}

	if p.privateKey != y.privateKey {
		return false
	}

	return reflect.DeepEqual(p.options, y.options)
}
