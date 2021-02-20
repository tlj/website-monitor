package monitors

import (
	"fmt"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"time"
	"website-monitor/result"

	"github.com/go-rod/rod"
)

type HttpRenderMonitor struct {
	renderServer string
}

func NewHttpRenderMonitor(renderServer string) *HttpRenderMonitor {
	return &HttpRenderMonitor{
		renderServer: renderServer,
	}
}

func (jm *HttpRenderMonitor) Check(check Check) (*result.Results, error) {
	l, err := launcher.NewRemote(jm.renderServer)
	if err != nil {
		return nil, fmt.Errorf("error connecting to rod at %s: %s", jm.renderServer, err)
	}
	l.Set("window-size", "1920,1080")

	r := rod.New().Client(l.Client()).Timeout(10 * time.Second)
	if err := r.Connect(); err != nil {
		return nil, err
	}

	p, err := r.Page(proto.TargetCreateTarget{URL: check.Url})
	if err != nil {
		return nil, err
	}

	if err = p.Timeout(3 * time.Second).WaitLoad(); err != nil {
		return nil, err
	}

	results := &result.Results{}
	for _, contentCheck := range check.ContentChecks {
		res, err := contentCheck.CheckRender(p)
		results.Results = append(results.Results, result.Result{
			ContentChecker: contentCheck,
			Result:         res,
			Err:            err,
		})
	}

	return results, nil
}

func (jm *HttpRenderMonitor) Type() string {
	return "HttpRenderMonitor"
}

func (jm *HttpRenderMonitor) Equal(y *HttpRenderMonitor) bool {
	return jm.renderServer == y.renderServer
}