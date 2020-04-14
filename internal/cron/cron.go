package cron

import (
	"strings"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

type cron struct {
	ActionName string
	Action     func()
	Delay      time.Duration
}

func (c *cron) wrapAction(action func()) func() {
	return func() {
		now := time.Now().UTC()
		log.Infof("Start %sing stage", c.ActionName)
		action()
		log.Infof("Finish %sing stage", c.ActionName)
		log.Infof("%sing stage elapsed %s", strings.Title(c.ActionName), time.Now().UTC().Sub(now).String())
		if err := db.DbMgr.SetLastActionDate(c.ActionName, now); err != nil {
			log.Errorf("can't save last_action date for %s: %v", c.ActionName, err)
		}
	}
}

func (c *cron) Run() {
	// override user action with our wrapper with extra logic
	c.Action = c.wrapAction(c.Action)
	now := time.Now().UTC()
	next := time.Now().UTC().Add(c.Delay)
	// start task immediately if delay if over against provided delay
	if now.After(next) {
		log.Infof("Starting %ving stage immediately, because provided delay is over after last fetch: %v | configured delay is: %v",
			c.ActionName, time.Now().UTC().Sub(now).String(), c.Delay)
		c.Action()
	}

	for range time.NewTicker(c.Delay).C {
		c.Action()
	}
}

func Run(actionName string, delay time.Duration, action func()) {
	c := cron{Action: action, ActionName: actionName, Delay: delay}
	c.Run()
}
