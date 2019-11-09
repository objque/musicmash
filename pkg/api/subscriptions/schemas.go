package subscriptions

type Subscription struct {
	ArtistID int64 `json:"artist_id" gorm:"unique_index:idx_user_name_artist_id"`
}
