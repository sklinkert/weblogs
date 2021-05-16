package main

import (
	"bufio"
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/sklinkert/weblogs/pkg/request"
	"os"
	"regexp"
	"strconv"
	"time"
)

var conf string
var format string
var logFile string

func init() {
	flag.StringVar(&format, "format", `$remote_addr [$time_local] "$request" $status $request_length $body_bytes_sent $request_time "$t_size" $read_time $gen_time`, "Log format")
	flag.StringVar(&logFile, "log", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
}

var topPath = map[string]int{}
var uniqueVisitors = map[uint32]int{}

func main() {
	flag.Parse()

	var err error
	file, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	re := regexp.MustCompile(`([\d\.]+) \- \- \[(\d\d\/\w+\/\d\d\d\d:\d\d:\d\d:\d\d \+\d{4})\] "\w+ (.*) HTTP/\d\.\d" (\d+) \d+ "(.*)" "(.*)"`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)
		//log.Info(matches)
		if len(matches) == 0 {
			log.Info("no match: ", line)
			continue
		}

		var (
			remoteAddr    = matches[0][1]
			localTimeStr  = matches[0][2]
			path          = matches[0][3]
			statusCodeStr = matches[0][4]
			referrer      = matches[0][5]
			userAgent     = matches[0][6]
		)

		statusCode, err := strconv.ParseInt(statusCodeStr, 10, 32)
		if err != nil {
			log.WithError(err).Warnf("Cannot parse status code: %q in line %s", statusCodeStr, line)
			continue
		}

		localTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", localTimeStr)
		if err != nil {
			log.WithError(err).Warnf("Cannot parse time: %q in line %s", localTimeStr, line)
			continue
		}

		var req = request.New(localTime, int(statusCode), remoteAddr, path, referrer, userAgent)

		//log.WithFields(log.Fields{
		//	"RemoteAddr": req.RemoteAddr,
		//	"LocalTime": req.LocalTime,
		//	"Path": req.Path,
		//	"Status": req.StatusCode,
		//	"Referrer": req.Referrer,
		//	"UserAgent": req.UserAgent,
		//	"FingerPrint": req.FingerPrint(),
		//	"IsBot": req.IsBot(),
		//	"Time": req.LocalTime,
		//}).Info("Matched")

		topPath[req.Path]++
		uniqueVisitors[req.FingerPrint()]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Infof("Visitors: %d", len(uniqueVisitors))
}
