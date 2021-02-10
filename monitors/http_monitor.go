package monitors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpMonitor struct{}

func (jm *HttpMonitor) Check(check Check) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, check.Url, nil)
	if err != nil {
		return false, err
	}
	for k, v := range check.Headers {
		req.Header.Add(k, v)
	}

	hc := http.Client{}
	hc.Timeout = 5 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != check.ExpectedStatusCode {
		return false, fmt.Errorf("invalid statuscode: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	for _, contentCheck := range check.ContentChecks {
		result, err := contentCheck.Check(ioutil.NopCloser(bytes.NewBuffer(body)))
		if !result {
			return false, err
		}
	}

	return true, nil
}
