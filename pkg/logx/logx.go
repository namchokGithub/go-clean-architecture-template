package logx

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once sync.Once
	log  *logrus.Entry
)

// Init configures the global logger. Safe to call multiple times; only first call takes effect.
func Init(level string) {
	once.Do(func() {
		l := logrus.New()
		l.SetOutput(os.Stdout)
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			lvl = logrus.InfoLevel
		}
		l.SetLevel(lvl)
		log = logrus.NewEntry(l)
	})
}

// GetLog returns the global logrus entry, initializing with "info" level if not yet done.
func GetLog() *logrus.Entry {
	if log == nil {
		Init("info")
	}
	return log
}
