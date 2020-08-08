package db

import (
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/stretchr/testify/assert"
)

func (t *testDBSuite) TestSubscriptions_Create() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))

	// action
	err := Mgr.CreateSubscription(&Subscription{
		UserName: vars.UserObjque,
		ArtistID: 1,
	})

	// assert
	assert.NoError(t.T(), err)
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, nil)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, subs[0].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_Get() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSPY,
		Poster: vars.PosterSimple,
	}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1, 2}))

	// action
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 2)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, subs[0].ArtistPoster)

	assert.Equal(t.T(), uint64(0), subs[1].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[1].UserName)
	assert.Equal(t.T(), int64(2), subs[1].ArtistID)
	assert.Equal(t.T(), vars.ArtistSPY, subs[1].ArtistName)
	assert.Equal(t.T(), vars.PosterSimple, subs[1].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_Get_Limit() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSPY,
		Poster: vars.PosterSimple,
	}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1, 2}))

	// action
	var limit uint64 = 1
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, &GetSubscriptionsOpts{Limit: &limit})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, subs[0].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_Get_Offset() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSPY,
		Poster: vars.PosterSimple,
	}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1, 2}))

	// action
	var offset uint64 = 1
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, &GetSubscriptionsOpts{Offset: &offset})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(2), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSPY, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterSimple, subs[0].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_Get_Limit_Offset() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSPY,
		Poster: vars.PosterSimple,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistArchitects,
		Poster: vars.PosterGiant,
	}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1, 2, 3}))

	// action
	var (
		limit  uint64 = 1
		offset uint64 = 1
	)
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, &GetSubscriptionsOpts{
		Limit:  &limit,
		Offset: &offset,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 1)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(2), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistSPY, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterSimple, subs[0].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_Get_SortType() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSkrillex,
		Poster: vars.PosterMiddle,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistSPY,
		Poster: vars.PosterSimple,
	}))
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{
		Name:   vars.ArtistArchitects,
		Poster: vars.PosterGiant,
	}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1, 2, 3}))

	// action
	const sortType = "desc"
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, &GetSubscriptionsOpts{SortType: sortType})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 3)
	// id always equals to zero, because it doesn't load in query
	assert.Equal(t.T(), uint64(0), subs[0].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(3), subs[0].ArtistID)
	assert.Equal(t.T(), vars.ArtistArchitects, subs[0].ArtistName)
	assert.Equal(t.T(), vars.PosterGiant, subs[0].ArtistPoster)

	assert.Equal(t.T(), uint64(0), subs[1].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[1].UserName)
	assert.Equal(t.T(), int64(2), subs[1].ArtistID)
	assert.Equal(t.T(), vars.ArtistSPY, subs[1].ArtistName)
	assert.Equal(t.T(), vars.PosterSimple, subs[1].ArtistPoster)

	assert.Equal(t.T(), uint64(0), subs[2].ID)
	assert.Equal(t.T(), vars.UserObjque, subs[2].UserName)
	assert.Equal(t.T(), int64(1), subs[2].ArtistID)
	assert.Equal(t.T(), vars.ArtistSkrillex, subs[2].ArtistName)
	assert.Equal(t.T(), vars.PosterMiddle, subs[2].ArtistPoster)
}

func (t *testDBSuite) TestSubscriptions_Get_Empty() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1}))

	// action
	subs, err := Mgr.GetUserSubscriptions(vars.UserBot, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}

func (t *testDBSuite) TestSubscriptions_UnSubscribe() {
	// arrange
	assert.NoError(t.T(), Mgr.EnsureArtistExists(&Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), Mgr.SubscribeUser(vars.UserObjque, []int64{1}))
	subs, err := Mgr.GetUserSubscriptions(vars.UserObjque, nil)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), vars.UserObjque, subs[0].UserName)
	assert.Equal(t.T(), int64(1), subs[0].ArtistID)

	// action
	err = Mgr.UnSubscribeUser(vars.UserObjque, []int64{1})

	// assert
	assert.NoError(t.T(), err)
	subs, err = Mgr.GetUserSubscriptions(vars.UserObjque, nil)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), subs, 0)
}
