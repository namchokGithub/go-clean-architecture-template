package protocol

import (
	"github.com/robfig/cron/v3"
	"gitlab.company.com/projectname/pkg/logx"
)

// StartCron registers and starts background cron jobs.
// No-ops when APP_ENV=local to prevent accidental job execution during development.
func (a *App) StartCron() {
	if a.conf.App.ENV == "local" {
		return
	}
	c := cron.New()
	log := logx.GetLog()

	// Register jobs:
	// c.AddFunc("@every 1m", func() { ... })

	c.Start()
	log.Info("cron started")
}
