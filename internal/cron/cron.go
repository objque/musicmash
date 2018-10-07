package cron

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/fetcher"
	"github.com/musicmash/musicmash/internal/log"
)

func isMustFetch() bool {
	last, err := db.DbMgr.GetLastActionDate(db.ActionFetch)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}

		log.Error(err)
		return false
	}

	diff := calcDiffHours(last.Date)
	log.Infof("LastAction fetch was at '%s'. Next fetch after %v hour",
		last.Date.Format("2006-01-02 15:04:05"),
		config.Config.Fetching.CountOfSkippedHoursToFetch-diff)
	return diff >= config.Config.Fetching.CountOfSkippedHoursToFetch
}

func Run() {
	for {
		if isMustFetch() {
			now := time.Now().UTC()
			log.Infof("Start fetching stage for '%s'...", now.String())
			fetcher.Fetch()
			log.Infof("Finish fetching stage '%s'...", time.Now().UTC().String())
			log.Infof("Elapsed time '%s'", time.Now().UTC().Sub(now).String())
			if err := db.DbMgr.SetLastActionDate(db.ActionFetch, now); err != nil {
				log.Error("can't save last_fetch date")
			}
		}

		time.Sleep(time.Hour)
	}
}
