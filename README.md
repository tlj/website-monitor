# Website Monitor

![example workflow](https://github.com/tlj/website-monitor/actions/workflows/test.yaml/badge.svg)
![example workflow](https://github.com/tlj/website-monitor/actions/workflows/main.yaml/badge.svg)
[![codecov](https://codecov.io/gh/tlj/website-monitor/branch/master/graph/badge.svg)](https://codecov.io/gh/tlj/website-monitor)
![Docker Image Version (latest semver)](https://img.shields.io/docker/v/thomaslandro/website-monitor)

Website monitor is a simple tool which will look for changes in
any HTTP hosted file (server side rendered html, json, etc).

## Usage

Add your configuration in config/config.yaml.

Dependency:
```shell
docker run -d --name website-renderer rodorg/rod:v0.91.1
```

Run with:
```shell
docker run -it --rm --links website-renderer -v $(pwd)/config.yaml:/app/config/config.yaml thomaslandro/website-monitor
```

docker-compose.yaml:
```yaml
  website-monitor:
    container_name: website-monitor
    image: thomaslandro/website-monitor:latest
    volumes:
      - ~/config/website-monitor:/app/config
    restart: always
    links:
      - website-renderer # optional if you want to use http_render monitor
    ports:
      - 2112:2112 # optional for prometheus metrics at /metrics
  
  # optional if you want to use http_render monitor
  website-renderer:
    container_name: website-renderer
    image: rodorg/rod:v0.91.1
    restart: always
```

## Configuration

Example:
```yaml
loglevel: info
defaults:
  type: "http"
  expected_status_code: 200 # http status code
  schedule: # optional
    interval: 60 # interval in seconds
    interval_variable_percentage: 20 # +/- 20% of the specified interval, making the range 48-72s
    days: "1-5" # every weekday (Mon-Fri)
    hours: "9-16" # between 9:00 and 16:59
  headers: # always send these headers in http requests
    User-Agent: "Mozilla/5.0"
  notifiers: # always send notifications on state change to these notifiers
    - name: Slack
      type: slack
      options:
        webhook: "https://hooks.slack.com/services/..."
    - name: PushSafer
      type: pushsafer
      options:
        private_key: "PRIVATE_KEY"
monitors:
  - name: "Monitored website feed"
    url: "https://www.monitored.website.example/feed.json"
    display_url: "https://www.monitored.website.example/"
    headers:
      Referer: "https://www.monitored.website.example/"
      Accept: "application/json"
    monitors:
      - name: SomeProperty
        type: json_path
        path: "//SomeProperty"
        value: "Whatever"
        is_expected: false
  - name: "Monitored website"
    url: "https://www.monitored.website.example/"
    type: http
    monitors:
      - name: SomeText
        type: regex
        value: "Some Text"
        is_expected: false
  - name: "Monitored website, two checks - one needed"
    url: "https://www.monitored.website.example/"
    require_some: true
    monitors:
      - name: SomeText
        type: regex
        value: "Some Text"
        is_expected: false        
      - name: SomeOtherText
        type: regex
        value: "Some Other Text"
        is_expected: false
  - name: "Don't use default schedule"
    url: "https://www.monitored.website.example/simple"
    monitors:
      - name: Some monitored text
        type: regex
        value: "Some monitored text"
        is_expected: true
    schedule:
      interval: 3600
      interval_variable_percentage: 0 
      days: ""
      hours: ""
  - name: "JS rendered website, with css selector"
    url: "https://www.monitored.website.example/js"
    type: http_render
    monitors:
      - name: Some text rendered only with JS
        type: HtmlRenderSelector
        path: "html body div#header h1#rendered"
        value: "A rendered header"
        is_expected: true
```