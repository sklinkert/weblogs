# weblogs
Web analytics from server logs

## Goals

- Use web server logs to get insights about your visitors and content
- No cookie consent nightmare
- Small tool. No requests to third-party servers

## Import nginx log to sqlite DB
```sh
go run cmd/weblogs/main.go -log /var/tmp/access.log
```
