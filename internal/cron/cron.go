package cron

import (
	"fmt"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

type cron struct {
	Action     func() error
	ActionName string
	Delay      time.Duration
}

func (c *cron) doActionAndUpdateLast() {
	// do action...
	log.Infof("calling %v action", c.ActionName)
	if err := c.Action(); err != nil {
		log.Errorf("%v action return err: %w", err)
		return
	}

	// update date when action was successful
	now := time.Now().UTC()
	if err := db.Mgr.SetLastActionDate(c.ActionName, now); err != nil {
		log.Errorf("can't save last_action date for %s: %v", c.ActionName, err)
		return
	}

	log.Infof("successfully update date for %v action", c.ActionName)
}

func (c *cron) Run() {
	// get last date when action was successful
	last, err := db.Mgr.GetLastActionDate(c.ActionName)
	if err != nil {
		log.Error(fmt.Errorf("tried to get last_action for %v stage: %w", c.ActionName, err))
		return
	}

	// check if action is outdated and we should start action now
	now := time.Now().UTC()
	previous := last.Date.Add(c.Delay)
	if now.After(previous) {
		log.Infof("%v action was too late, trigger it now", c.ActionName)
		c.doActionAndUpdateLast()
	}

	// schedule new ticker
	log.Infof("starting ticker with %v delay for %v action", c.Delay, c.ActionName)
	for range time.NewTicker(c.Delay).C {
		c.doActionAndUpdateLast()
	}
}

func Run(actionName string, delay time.Duration, action func() error) {
	scheduler := cron{
		Action:     action,
		ActionName: actionName,
		Delay:      delay,
	}
	scheduler.Run()
}
