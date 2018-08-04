package db

type Subscription struct {
	ID         int64  `gorm:"primary_key"`
	UserID     string `sql:"index"`
	ArtistName string
}

type SubscriptionMgr interface {
	IsUserSubscribedForArtist(userID, artistName string) bool
	FindAllUserSubscriptions(userID string) ([]*Subscription, error)
	EnsureSubscriptionExists(subscription *Subscription) error
}

func (mgr *AppDatabaseMgr) IsUserSubscribedForArtist(userID, artistName string) bool {
	subscription := Subscription{}
	if err := mgr.db.Where("user_id = ? and artist_name = ?", userID, artistName).First(&subscription).Error; err != nil {
		return false
	}

	return true
}

func (mgr *AppDatabaseMgr) FindAllUserSubscriptions(userID string) ([]*Subscription, error) {
	subscriptions := []*Subscription{}
	if err := mgr.db.Where("user_id = ?", userID).Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (mgr *AppDatabaseMgr) EnsureSubscriptionExists(subscription *Subscription) error {
	if !mgr.IsUserSubscribedForArtist(subscription.UserID, subscription.ArtistName) {
		return mgr.db.Create(subscription).Error
	}
	return nil
}
