package fetcher

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher/handlers/itunes"
	"github.com/objque/musicmash/internal/fetcher/v2"
	"github.com/objque/musicmash/internal/log"
)

func isMustFetch() bool {
	last, err := db.DbMgr.GetLastFetch()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}

		log.Error(err)
		return false
	}

	diff := calcDiffHours(last.Date)
	log.Infof("Last fetch was at '%s'. Next fetch after %v hour",
		last.Date.Format("2006-01-02 15:04:05"),
		config.Config.Fetching.CountOfSkippedHoursToFetch-diff)
	return diff >= config.Config.Fetching.CountOfSkippedHoursToFetch
}

func Run() {
	f := v2.Fetcher{}
	f.RegisterHandler(&itunes.AppleMusicHandler{})
	for {
		if isMustFetch() {
			now := time.Now().UTC()
			log.Infof("Start fetching stage for '%s'...", now.String())
			if err := f.FetchAndProcess(); err != nil {
				log.Error(err)
			} else {
				log.Infof("Finish fetching stage '%s'...", time.Now().UTC().String())
				db.DbMgr.SetLastFetch(now)
			}
			log.Infof("Elapsed time '%s'", time.Now().UTC().Sub(now).String())
		}

		time.Sleep(time.Hour)
	}
}
