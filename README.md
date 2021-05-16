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

## Statistics from DB

```sql
-- top OS
SELECT os, count(*) count FROM requests WHERE is_bot IS FALSE GROUP BY os ORDER by count DESC LIMIT 100;
Windows	20012
Android	11413
iOS	8370
macOS	7598
Linux	1652
```

```sql
-- Top referrer
SELECT referrer, count(*) count FROM requests WHERE is_bot IS FALSE GROUP BY referrer ORDER by count DESC LIMIT 100;
-	7017
https://www.google.com/	2495
https://www.topblogs.de	2331
https://etf.capital/inflation-etf/	1421
https://etf.capital/	1386
```

```sql
-- Requests by day
SELECT count(*), date(local_time) day FROM requests WHERE is_bot IS FALSE GROUP BY date(local_time) ORDER by day DESC LIMIT 100;
246114	2021-05-16
258788	2021-05-15
```

```sql
-- Top path
SELECT path, count(*) count FROM requests WHERE is_bot IS FALSE GROUP BY path ORDER by count DESC LIMIT 100;
/rss/	1705
/	1621
/favicon.ico	636
```

