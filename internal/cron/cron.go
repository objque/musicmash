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
			if err = db.DbMgr.SetLastActionDate(c.ActionName, time.Now().UTC()); err != nil {
				log.Errorf("can't save last_action date for '%s', do it manually", c.ActionName)
				return false
			}

			log.Infof("Last %s set as now(). Next in %v hour",
				c.ActionName, c.CountOfSkippedHoursToRun)
			return false
		}

		log.Error(err)
		return false
	}

	diff := calcDiffHours(last.Date)
	log.Infof("Last %s was at %s. Next in %v hour",
		c.ActionName,
		last.Date.Format("2006-01-02T15:04:05"),
		c.CountOfSkippedHoursToRun-diff)
	return diff >= c.CountOfSkippedHoursToRun
}

func Run(actionName string, countOfSkippedHours float64, action func()) {
	c := cron{Action: action, ActionName: actionName, CountOfSkippedHoursToRun: countOfSkippedHours}
	c.Run()
}
