package db

import (
	"time"
)

type User struct {
	CreatedAt time.Time
	ID        string `gorm:"primary_key"`
}

type UserMgr interface {
	CreateUser(user *User) error
	FindUserByID(id string) (*User, error)
	GetAllUsers() ([]*User, error)
	EnsureUserExists(userID string) error
}

func (mgr *AppDatabaseMgr) FindUserByID(id string) (*User, error) {
	user := User{}
	if err := mgr.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (mgr *AppDatabaseMgr) GetAllUsers() ([]*User, error) {
	var users = make([]*User, 0)
	return users, mgr.db.Find(&users).Error
}

func (mgr *AppDatabaseMgr) CreateUser(user *User) error {
	return mgr.db.Create(user).Error
}

func (mgr *AppDatabaseMgr) EnsureUserExists(userID string) error {
	_, err := mgr.FindUserByID(userID)
	if err != nil {
		return mgr.CreateUser(&User{ID: userID})
	}
	return nil
}
