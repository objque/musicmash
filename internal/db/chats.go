package db

import (
	"time"
)

type Chat struct {
	CreatedAt time.Time
	ID        int64  `gorm:"primary_key"`
	UserID    string `sql:"index"`
}

type ChatMgr interface {
	FindChatByUserID(userID string) (*int64, error)
	EnsureChatExists(chat *Chat) error
}

func (mgr *AppDatabaseMgr) FindChatByUserID(userID string) (*int64, error) {
	chat := Chat{}
	if err := mgr.db.Where("user_id = ?", userID).First(&chat).Error; err != nil {
		return nil, err
	}

	return &chat.ID, nil
}

func (mgr *AppDatabaseMgr) EnsureChatExists(chat *Chat) error {
	_, err := mgr.FindChatByUserID(chat.UserID)
	if err != nil {
		return mgr.db.Create(chat).Error
	}
	return nil
}
