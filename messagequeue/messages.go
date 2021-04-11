package messagequeue

type ScheduleJob struct {
	MonitorID int64  `json:"monitor_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type CrawlCheckerResult struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Result  bool   `json:"result"`
	Err     error  `json:"err"`
}

type CrawlResult struct {
	MonitorID int64                `json:"monitor_id"`
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	Results   []CrawlCheckerResult `json:"results"`
	Result    bool                 `json:"result"`
}
