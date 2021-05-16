package request

import (
	"github.com/mileusna/useragent"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hash/fnv"
	"time"
)

var db *gorm.DB

type Request struct {
	RemoteAddr string    `gorm:"primaryKey"`
	LocalTime  time.Time `gorm:"primaryKey"`
	Path       string    `gorm:"primaryKey"`
	StatusCode int
	Referrer   string
	UserAgent  string
	IsBot      bool
	IsMobile   bool
	IsDesktop  bool
	IsTablet   bool
	OS         string
	OSVersion  string
	Device     string
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (req *Request) FingerPrint() uint32 {
	return hash(req.RemoteAddr + req.UserAgent)
}

func New(localTime time.Time, statusCode int, remoteAddr, path, referrer, userAgent string) *Request {
	parsedUA := ua.Parse(userAgent)
	return &Request{
		RemoteAddr: remoteAddr,
		LocalTime:  localTime,
		Path:       path,
		StatusCode: statusCode,
		Referrer:   referrer,
		UserAgent:  userAgent,
		IsBot:      parsedUA.Bot,
		IsDesktop:  parsedUA.Desktop,
		IsTablet:   parsedUA.Tablet,
		OS:         parsedUA.OS,
		OSVersion:  parsedUA.OSVersion,
		Device:     parsedUA.Device,
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
