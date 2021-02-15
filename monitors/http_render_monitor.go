package monitors

import (
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"time"
	"website-monitor/result"

	"github.com/go-rod/rod"
)

type HttpRenderMonitor struct{}

func (jm *HttpRenderMonitor) Check(check Check) (*result.Results, error) {
	l, err := launcher.NewRemote("ws://localhost:9222")
	if err != nil {
		return nil, err
	}
	l.Set("window-size", "1920,1080")

	r := rod.New().Client(l.Client()).Timeout(5 * time.Second)
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
