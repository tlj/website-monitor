package app

import (
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"website-monitor/monitors"
)

type Config struct {
	LogLevel string              `yaml:"loglevel"`
	Default  *monitors.Monitor   `yaml:"defaults"`
	Monitors []*monitors.Monitor `yaml:"monitors"`
}

func (c *Config) LoadConfigFromFile(filename string) error {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return c.LoadConfig(configData)
}

func (c *Config) LoadConfig(configData []byte) error {
	err := yaml.Unmarshal(configData, c)
	if err != nil {
		return err
	}

	for _, chk := range c.Monitors {
		if chk.DisplayUrl == "" {
			chk.DisplayUrl = chk.Url
		}
		if chk.Headers == nil {
			chk.Headers = make(map[string]string)
		}
		if _, ok := chk.Headers["Referer"]; !ok {
			chk.Headers["Referer"] = chk.DisplayUrl
		}
		if c.Default != nil {
			if c.Default.Notifiers != nil && chk.Notifiers == nil {
				chk.Notifiers = c.Default.Notifiers
			}
			if c.Default.Scheduler != nil && chk.Scheduler == nil {
				chk.Scheduler = c.Default.Scheduler
			}
			if c.Default.Type != "" && chk.Type == "" {
				chk.Type = c.Default.Type
			}
			if c.Default.ExpectedStatusCode != 0 && chk.ExpectedStatusCode == 0 {
				chk.ExpectedStatusCode = c.Default.ExpectedStatusCode
			}

			for k, v := range c.Default.Headers {
				if _, ok := chk.Headers[k]; !ok {
					chk.Headers[k] = v
				}
			}

		}
	}

	return nil
}

/*
func (c *Config) Parse(data []byte) error {
	cfg := &configFile{}
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	c.LogLevel = cfg.LogLevel

	for k, mCheck := range cfg.Monitors {
		if mCheck.Url == "" && mCheck.DisplayUrl == "" {
			return fmt.Errorf("monitor %d (%s) is missing url and display_url", k, mCheck.Name)
		}
		if mCheck.Url == "" {
			mCheck.Url = mCheck.DisplayUrl
		}
		if mCheck.DisplayUrl == "" {
			mCheck.DisplayUrl = mCheck.Url
		}
		if mCheck.Type == "" {
			mCheck.Type = "http"
		}
		if mCheck.ExpectedStatusCode == 0 {
			mCheck.ExpectedStatusCode = cfg.Global.ExpectedStatusCode
		}
		if mCheck.ExpectedStatusCode == 0 {
			return fmt.Errorf("monitor %d (%s) has invalid required status code: %d", k, mCheck.Name, mCheck.ExpectedStatusCode)
		}
		check := monitors.Monitor{
			Name:               mCheck.Name,
			Url:                mCheck.Url,
			DisplayUrl:         mCheck.DisplayUrl,
			RenderServerURN:    mCheck.RenderServerURN,
			Type:               mCheck.Type,
			Headers:            mCheck.Headers,
			ExpectedStatusCode: mCheck.ExpectedStatusCode,
			ContentChecks:      nil,
			Notifiers:          nil,
		}
		if mCheck.Interval < 1 {
			mCheck.Interval = cfg.Global.Interval
		}
		if mCheck.Interval < 1 {
			return fmt.Errorf("monitor %d (%s) has invalid interval: %d", k, mCheck.Name, mCheck.Interval)
		}
		if mCheck.IntervalVariablePercentage == nil {
			mCheck.IntervalVariablePercentage = cfg.Global.IntervalVariablePercentage
		}

		if mCheck.RegexExpected != "" {
			check.ContentChecks = append(check.ContentChecks, content_checkers.NewRegexChecker(mCheck.RegexExpected, mCheck.RegexExpected, true))
		}

		if mCheck.RegexNotExpected != "" {
			check.ContentChecks = append(check.ContentChecks, content_checkers.NewRegexChecker(mCheck.RegexNotExpected, mCheck.RegexNotExpected, false))
		}

		if mCheck.RenderServerURN == "" {
			mCheck.RenderServerURN = cfg.Global.RenderServerURN
		}
		check.RenderServerURN = mCheck.RenderServerURN

		check.Headers = make(map[string]string)
		for k, v := range cfg.Global.Headers {
			check.Headers[k] = v
		}
		for k, v := range mCheck.Headers {
			check.Headers[k] = v
		}

		if _, ok := check.Headers["Referer"]; !ok {
			if check.DisplayUrl != "" && check.Url != "" {
				check.Headers["Referer"] = check.DisplayUrl
			} else {
				u, _ := url.Parse(check.Url)
				check.Headers["Referer"] = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
			}
		}

		if mCheck.Schedule == nil {
			mCheck.Schedule = cfg.Global.Schedule
		}

		var hours []int
		var days []time.Weekday

		if mCheck.Schedule != nil {
			if mCheck.Schedule.Hours != "" {
				startEndHours := strings.Split(mCheck.Schedule.Hours, "-")
				startHour, _ := strconv.Atoi(startEndHours[0])
				endHour, _ := strconv.Atoi(startEndHours[1])
				for i := startHour; i <= endHour; i++ {
					hours = append(hours, i)
				}
			}

			if mCheck.Schedule.Days != "" {
				startEndDays := strings.Split(mCheck.Schedule.Days, "-")
				startDay, _ := strconv.Atoi(startEndDays[0])
				endDay, _ := strconv.Atoi(startEndDays[1])
				for i := startDay; i <= endDay; i++ {
					days = append(days, time.Weekday(i))
				}
			}

			check.Scheduler = scheduler.NewScheduler(time.Duration(mCheck.Interval)*time.Second, mCheck.IntervalVariablePercentage, hours, days)
		}

		for _, n := range cfg.Global.NotifiersConfig {
			var notifier notifiers.Notifier
			switch n["type"] {
			case "slack":
				notifier = notifiers.NewSlackNotifier(n["webhook"])
			default:
				return fmt.Errorf("unsupported notifiers config: %s", n["type"])
			}
			check.Notifiers = append(check.Notifiers, notifier)
		}

		for _, n := range mCheck.NotifiersConfig {
			var notifier notifiers.Notifier
			switch n["type"] {
			case "slack":
				notifier = notifiers.NewSlackNotifier(n["webhook"])
			default:
				return fmt.Errorf("unsupported notifiers config: %s", n["type"])
			}
			check.Notifiers = append(check.Notifiers, notifier)
		}

		for _, cc := range mCheck.ContentChecksConfig {
			var expected string
			var expectedToExist bool
			if v, ok := cc["expected"]; ok {
				expected = v
				expectedToExist = true
			}
			if v, ok := cc["not_expected"]; ok {
				expected = v
				expectedToExist = false
			}

			var contentCheck content_checkers.ContentChecker
			switch cc["type"] {
			case "JsonPath":
				contentCheck = content_checkers.NewJsonPathChecker(cc["name"], cc["path"], expected, expectedToExist)
			case "Regex":
				contentCheck = content_checkers.NewRegexChecker(cc["name"], expected, expectedToExist)
			case "HtmlXPath":
				contentCheck = content_checkers.NewHtmlXPathChecker(cc["name"], cc["path"], expected, expectedToExist)
			case "HtmlRenderSelector":
				contentCheck = content_checkers.NewHtmlRenderSelectorChecker(cc["name"], cc["path"], expected, expectedToExist)
			default:
				return fmt.Errorf("unsupported contentCheck config: %s", cc["type"])
			}

			check.ContentChecks = append(check.ContentChecks, contentCheck)
		}

		c.Monitors = append(c.Monitors, &check)
	}

	return nil
}
*/
