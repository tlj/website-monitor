package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
	"website-monitor/notifiers"
	"website-monitor/scheduler"
)

type configFile struct {
	LogLevel string `yaml:"loglevel"`
	Global   struct {
		Headers                    map[string]string   `yaml:"headers"`
		ExpectedStatusCode         int                 `yaml:"expected_status_code"`
		Interval                   int                 `yaml:"interval"`
		IntervalVariablePercentage *int                `yaml:"interval_variable_percentage"`
		NotifiersConfig            []map[string]string `yaml:"notifiers"`
		RenderServerURN            string              `yaml:"render_server_urn"`
		Schedule                   *struct {
			Days  string `yaml:"days"`
			Hours string `yaml:"hours"`
		} `yaml:"schedule"`
	} `yaml:"global"`
	Monitors []struct {
		Name                       string               `yaml:"name"`
		Url                        string               `yaml:"url"`
		DisplayUrl                 string               `yaml:"display_url"`
		RenderServerURN            string               `yaml:"render_server_urn"`
		Type                       monitors.MonitorType `yaml:"type"`
		Headers                    map[string]string    `yaml:"headers"`
		RegexNotExpected           string               `yaml:"regex_not_expected"`
		RegexExpected              string               `yaml:"regex_expected"`
		ExpectedStatusCode         int                  `yaml:"expected_status_code"`
		LastSeenState              bool                 `yaml:"last_seen_state"`
		ContentChecksConfig        []map[string]string  `yaml:"content_checks"`
		NotifiersConfig            []map[string]string  `yaml:"notifiers"`
		Interval                   int                  `yaml:"interval"`
		IntervalVariablePercentage *int                 `yaml:"interval_variable_percentage"`
		Schedule                   *struct {
			Days  string `yaml:"days"`
			Hours string `yaml:"hours"`
		} `yaml:"schedule"`
	} `yaml:"monitors"`
}

type Config struct {
	LogLevel string            `yaml:"loglevel"`
	Monitors []*monitors.Check `yaml:"monitors"`
}

func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = config.Parse(configData)
	if err != nil {
		return nil, err
	}

	return config, nil
}

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
		check := monitors.Check{
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
