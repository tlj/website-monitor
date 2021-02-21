# Website Monitor

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
global:
  expected_status_code: 200 # http status code
  interval: 60 # interval in seconds
  interval_variable_percentage: 20 # +/- 20% of the specified interval, making the range 48-72s
  schedule: # optional
    days: "1-5" # every weekday (Mon-Fri)
    hours: "9-16" # between 9:00 and 16:59
  headers: # always send these headers in http requests
    User-Agent: "Mozilla/5.0"
  notifiers: # always send notifications on state change to these notifiers
    - name: Slack
      type: slack
      webhook: "https://hooks.slack.com/services/..."
monitors:
  - name: "Monitored website feed"
    url: "https://www.monitored.website.example/feed.json"
    display_url: "https://www.monitored.website.example/"
    type: http
    headers:
      Referer: "https://www.monitored.website.example/"
      Accept: "application/json"
    content_checks:
      - name: SomeProperty
        type: JsonPath
        path: "//SomeProperty"
        not_expected: "Whatever"
  - name: "Monitored website feed"
    url: "https://www.monitored.website.example/"
    type: http
    content_checks:
      - name: SomeText
        type: Regex
        not_expected: "Some Text"
  - name: "Simpler config for website monitor"
    url: "https://www.monitored.website.example/simple"
    regex_expected: "Some monitored text" # short-hand form for the full content_check of "Regex" type
  - name: "Don't use default schedule"
    url: "https://www.monitored.website.example/simple"
    regex_expected: "Some monitored text"
    interval: 3600
    interval_variable_percentage: 0
    schedule: {}
  - name: "JS rendered website, with css selector"
    url: "https://www.monitored.website.example/js"
    type: http_render
    content_checks:
      - name: Some text rendered only with JS
        type: HtmlRenderSelector
        path: "html body div#header h1#rendered"
        expected: "A rendered header"
```