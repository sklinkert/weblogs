package request

import (
	"fmt"
	"github.com/mileusna/useragent"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hash/fnv"
	"strings"
	"time"
)

var db *gorm.DB

type Request struct {
	Fingerprint string    `gorm:"primaryKey"`
	LocalTime   time.Time `gorm:"primaryKey"`
	Path        string    `gorm:"primaryKey"`
	Method      string    `gorm:"primaryKey"`
	RemoteAddr  string
	StatusCode  int
	Referrer    string
	UserAgent   string
	IsBot       bool
	IsMobile    bool
	IsDesktop   bool
	IsTablet    bool
	OS          string
	OSVersion   string
	Device      string
}

func toFingerprint(remoteAddr, userAgent string) string {
	key := strings.Join([]string{remoteAddr, userAgent}, ":")
	h := fnv.New64a()
	h.Write([]byte(key))
	return fmt.Sprintf("%d", h.Sum64())
}

func New(localTime time.Time, statusCode int, method, remoteAddr, path, referrer, userAgent string) *Request {
	parsedUA := ua.Parse(userAgent)

	return &Request{
		RemoteAddr:  remoteAddr,
		LocalTime:   localTime,
		Method:      method,
		Path:        path,
		StatusCode:  statusCode,
		Referrer:    referrer,
		UserAgent:   userAgent,
		IsBot:       parsedUA.Bot,
		IsDesktop:   parsedUA.Desktop,
		IsTablet:    parsedUA.Tablet,
		OS:          parsedUA.OS,
		OSVersion:   parsedUA.OSVersion,
		Device:      parsedUA.Device,
		Fingerprint: toFingerprint(remoteAddr, userAgent),
	}
}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("weblogs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Request{})
}

func (req *Request) Save() error {
	return db.Create(&req).Error
}
