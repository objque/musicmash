package artists

import (
	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/itunes"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

func EnsureExists(artists []string) (found []string, notFound []string) {
	found, notFound = []string{}, []string{}
	for _, userArtist := range artists {
		dbArtist, err := db.DbMgr.FindArtistByName(userArtist)
		// artist already exists
		if err == nil {
			// override artist name that was provided by a user
			// because user may send 'skriLLLex', but we store 'Skrillex'
			found = append(found, dbArtist.Name)
			continue
		}
		// another db err raised
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(errors.Wrapf(err, "tried to get artist '%s' from the db", userArtist))
			continue
		}

		artist, err := itunes.FindArtistID(userArtist)
		if err != nil {
			if err == itunes.ArtistNotFoundErr {
				notFound = append(notFound, userArtist)
				err = errors.Wrap(err, userArtist)
			}

			log.Error(err)
			continue
		}

		err = db.DbMgr.CreateArtist(&db.Artist{Name: artist.Name, StoreID: artist.StoreID})
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to add new artist '%s'", userArtist))
		}
		found = append(found, artist.Name)
	}
	return
}
