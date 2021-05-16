package request

import (
	ua "github.com/mileusna/useragent"
	"hash/fnv"
	"time"
)

type Request struct {
	RemoteAddr      string
	LocalTime       time.Time
	Path            string
	StatusCode      int
	Referrer        string
	UserAgent       string
	parsedUserAgent ua.UserAgent
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (req *Request) FingerPrint() uint32 {
	return hash(req.RemoteAddr + req.UserAgent)
}

func (req *Request) IsBot() bool {
	return req.parsedUserAgent.Bot
}

func New(localTime time.Time, statusCode int, remoteAddr, path, referrer, userAgent string) *Request {
	return &Request{
		RemoteAddr:      remoteAddr,
		LocalTime:       localTime,
		Path:            path,
		StatusCode:      statusCode,
		Referrer:        referrer,
		UserAgent:       userAgent,
		parsedUserAgent: ua.Parse(userAgent),
	}
}
