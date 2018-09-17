package fetcher

import (
	"time"

	"github.com/jinzhu/gorm"
	itunesProvider "github.com/objque/musicmash/internal/clients/itunes"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher/handlers/itunes"
	"github.com/objque/musicmash/internal/fetcher/handlers/yandex"
	"github.com/objque/musicmash/internal/fetcher/v2"
	"github.com/objque/musicmash/internal/log"
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
	provider := itunesProvider.NewProvider(config.Config.Store.URL, config.Config.Store.Token)
	f := v2.Fetcher{Provider: provider}
	f.RegisterHandler(&itunes.AppleMusicHandler{})
	f.RegisterHandler(yandex.New("https://music.yandex.ru", provider))
	for {
		if isMustFetch() {
			now := time.Now().UTC()
			log.Infof("Start fetching stage for '%s'...", now.String())
			if err := f.FetchAndProcess(); err != nil {
				log.Error(err)
			} else {
				// NOTE (m.kalinin): release.created_at must be after last_fetch, because
				// notifier uses that time to find new releases;
				log.Infof("Finish fetching stage '%s'...", time.Now().UTC().String())
				db.DbMgr.SetLastActionDate(db.ActionFetch, now)
			}
			log.Infof("Elapsed time '%s'", time.Now().UTC().Sub(now).String())
		}

		time.Sleep(time.Hour)
	}
}
