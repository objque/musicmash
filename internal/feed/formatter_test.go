package feed

import (
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFormatter_ToRss(t *testing.T) {
	// arrange
	const title = "Fresh releases from musicmash.me for objque@me"
	const link = "https://musicmash.me/news.xml"
	const desc = "Fresh releases from your artists"
	formatter := Formatter{Title: title, Link: link, Description: desc}
	releases := []*db.Release{
		{StoreName: testutil.StoreDeezer, Title: testutil.ReleaseArchitectsHollyHell, ArtistName: testutil.ArtistArchitects},
		{StoreName: testutil.StoreApple, Title: testutil.ReleaseSkrillexRecess, ArtistName: testutil.ArtistSkrillex},
	}

	// action
	rss := formatter.ToRawRSS(releases)

	// assert
	assert.Equal(t, link, rss.Link.Href)
	assert.Equal(t, title, rss.Title)
	assert.Equal(t, desc, rss.Description)
	assert.Len(t, rss.Items, 2)
	assert.Equal(t, makeRssDescription(releases[0]), rss.Items[0].Title)
	assert.Equal(t, makeRssDescription(releases[0]), rss.Items[0].Description)
	assert.Equal(t, makeRssDescription(releases[1]), rss.Items[1].Title)
	assert.Equal(t, makeRssDescription(releases[1]), rss.Items[1].Description)
	// NOTE (m.kalinin): replace hardcoded urls after extract into config
	assert.Equal(t, "https://www.deezer.com/en/album/", rss.Items[0].Link.Href)
	assert.Equal(t, "https://itunes.apple.com/us/album/", rss.Items[1].Link.Href)
}
