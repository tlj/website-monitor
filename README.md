# Website Monitor

Website monitor is a simple tool which will look for changes in
any HTTP hosted file (server side rendered html, json, etc).

## Usage

Add your configuration in config/config.yaml.

Run with:
```shell
docker run -v $(pwd)/config.yaml:/app/config/config.yaml thomaslandro/website-monitor
```

Example:
```yaml
global:
  expected_status_code: 200
  interval: 30
  headers:
    User-Agent: "Mozilla/5.0"
  notifiers:
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
        expected: "Whatever"
        expected_to_exist: false
  - name: "Monitored website feed"
    url: "https://www.monitored.website.example/"
    type: http
    content_checks:
      - name: SomeText
        type: Regex
        regex: "Some Text"
        expected_to_exist: false
```