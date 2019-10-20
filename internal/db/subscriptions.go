package db

type Subscription struct {
	ID       uint64 `json:"-"         gorm:"primary_key"     sql:"AUTO_INCREMENT"`
	UserName string `json:"user_name" gorm:"unique_index:idx_user_name_artist_id"`
	ArtistID int64  `json:"artist_id" gorm:"unique_index:idx_user_name_artist_id"`
}

type SubscriptionMgr interface {
	GetSimpleUserSubscriptions(userName string) ([]int64, error)
	GetUserSubscriptions(userName string) ([]*Subscription, error)
	GetArtistsSubscriptions(artists []int64) ([]*Subscription, error)
	SubscribeUser(userName string, artists []int64) error
	UnSubscribeUser(userName string, artists []int64) error
}

func (mgr *AppDatabaseMgr) GetUserSubscriptions(userName string) ([]*Subscription, error) {
	subs := []*Subscription{}
	err := mgr.db.Where("user_name = ?", userName).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (mgr *AppDatabaseMgr) GetArtistsSubscriptions(artists []int64) ([]*Subscription, error) {
	subs := []*Subscription{}
	err := mgr.db.Where("artist_id in (?)", artists).Order("user_name").Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (mgr *AppDatabaseMgr) GetSimpleUserSubscriptions(userName string) ([]int64, error) {
	ids := []int64{}
	err := mgr.db.Table("subscriptions").Where("user_name = ?", userName).Pluck("artist_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (mgr *AppDatabaseMgr) SubscribeUser(userName string, artists []int64) error {
	var q = "insert ignore into subscriptions (user_name, artist_id) values (?, ?)"
	if mgr.GetDialectName() == "sqlite3" {
		q = "insert or ignore into subscriptions (user_name, artist_id) values (?, ?)"
	}

	query, err := mgr.db.DB().Prepare(q)
	if err != nil {
		return err
	}

	for _, artistID := range artists {
		_, err := query.Exec(userName, artistID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mgr *AppDatabaseMgr) UnSubscribeUser(userName string, artists []int64) error {
	const query = "delete from subscriptions where user_name = ? and artist_id in (?)"
	return mgr.db.Exec(query, userName, artists).Error
}
