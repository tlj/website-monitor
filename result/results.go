package result

import "website-monitor/content_checkers"

type Results struct {
	Results []Result
}

func (r *Results) AllTrue() bool {
	for _, result := range r.Results {
		if !result.Result {
			return false
		}
	}

	return true
}

func (r *Results) SomeTrue() bool {
	for _, result := range r.Results {
		if result.Result {
			return true
		}
	}

	return false
}

type Result struct {
	ContentChecker content_checkers.ContentChecker
	Result         bool
	Err            error
}

