package monitors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"website-monitor/result"
)

type HttpMonitor struct{}

func (jm *HttpMonitor) Check(check Monitor) (*result.Results, error) {
	req, err := http.NewRequest(http.MethodGet, check.Url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range check.Headers {
		req.Header.Add(k, v)
	}

	hc := http.Client{}
	hc.Timeout = 5 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != check.ExpectedStatusCode {
		return nil, fmt.Errorf("invalid statuscode: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := &result.Results{}
	for _, contentCheck := range check.ContentChecks {
		res, err := contentCheck.Check(ioutil.NopCloser(bytes.NewBuffer(body)))
		results.Results = append(results.Results, result.Result{
			ContentChecker: contentCheck,
			Result: res,
			Err: err,
		})
	}

	return results, nil
}

func (jm *HttpMonitor) Type() string {
	return "HttpMonitor"
}
