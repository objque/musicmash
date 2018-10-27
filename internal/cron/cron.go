package cron

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

type cron struct {
	ActionName               string
	Action                   func()
	CountOfSkippedHoursToRun float64
}

func (c *cron) Run() {
	for {
		if !c.IsMustFetch() {
			time.Sleep(time.Minute * 15)
			continue
		}

		now := time.Now().UTC()
		log.Infof("Start %sing stage for '%s'...", c.ActionName, now.String())
		c.Action()
		log.Infof("Finish %sing stage '%s'...", c.ActionName, time.Now().UTC().String())
		log.Infof("Elapsed time '%s' for %s", time.Now().UTC().Sub(now).String(), c.ActionName)
		if err := db.DbMgr.SetLastActionDate(c.ActionName, now); err != nil {
			log.Errorf("can't save last_action date for '%s'", c.ActionName)
		}
	}
}

func (c *cron) IsMustFetch() bool {
	last, err := db.DbMgr.GetLastActionDate(c.ActionName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("LastAction for '%s' not found, will start now...", c.ActionName)
			return true
		}

		log.Error(err)
		return false
	}

	diff := calcDiffHours(last.Date)
	log.Infof("LastAction '%s' was at '%s'. Next fetch after %v hour",
		c.ActionName,
		last.Date.Format("2006-01-02 15:04:05"),
		c.CountOfSkippedHoursToRun-diff)
	return diff >= c.CountOfSkippedHoursToRun
}

func Run(actionName string, countOfSkippedHours float64, action func()) {
	c := cron{Action: action, ActionName: actionName, CountOfSkippedHoursToRun: countOfSkippedHours}
	c.Run()
}
