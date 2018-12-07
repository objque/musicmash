package db

import "time"

type User struct {
	CreatedAt time.Time
	Name      string `gorm:"primary_key"`
}

type UserMgr interface {
	CreateUser(user *User) error
	FindUserByName(name string) (*User, error)
	GetAllUsers() ([]*User, error)
	EnsureUserExists(userID string) error
	GetUsersWithReleases(date time.Time) ([]string, error)
}

func (mgr *AppDatabaseMgr) FindUserByName(id string) (*User, error) {
	user := User{}
	if err := mgr.db.Where("name = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (mgr *AppDatabaseMgr) GetAllUsers() ([]*User, error) {
	var users = []*User{}
	return users, mgr.db.Find(&users).Error
}

func (mgr *AppDatabaseMgr) CreateUser(user *User) error {
	return mgr.db.Create(user).Error
}

func (mgr *AppDatabaseMgr) EnsureUserExists(name string) error {
	_, err := mgr.FindUserByName(name)
	if err != nil {
		return mgr.CreateUser(&User{Name: name})
	}
	return nil
}

func (mgr *AppDatabaseMgr) GetUsersWithReleases(date time.Time) ([]string, error) {
	// Returns list of users that subscribed for artists that released/announced a new release
	const query = "select user_name from subscriptions where artist_name in (select artist_name from releases where created_at >= ?) group by user_name"
	users := []string{}
	if err := mgr.db.Raw(query, date).Pluck("user_name", &users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
