package db

import "errors"

type NotificationSettings struct {
	ID       int    `json:"-"      gorm:"primary_key" sql:"AUTO_INCREMENT"`
	UserName string `json:"-"`
	Service  string `json:"service"`
	Data     string `json:"data"`
}

type NotificationSettingsMgr interface {
	FindNotificationSettings(userName string) ([]*NotificationSettings, error)
	EnsureNotificationSettingsExists(settings *NotificationSettings) error
	UpdateNotificationSettings(settings *NotificationSettings) error
}

func (mgr *AppDatabaseMgr) FindNotificationSettings(userName string) ([]*NotificationSettings, error) {
	settings := []*NotificationSettings{}
	if err := mgr.db.Where("user_name = ?", userName).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (mgr *AppDatabaseMgr) EnsureNotificationSettingsExists(settings *NotificationSettings) error {
	userSettings, err := mgr.FindNotificationSettings(settings.UserName)
	if err != nil {
		return err
	}
	if len(userSettings) == 0 {
		return mgr.db.Create(settings).Error
	}
	return nil
}

func (mgr *AppDatabaseMgr) UpdateNotificationSettings(settings *NotificationSettings) error {
	existing := []*NotificationSettings{}
	err := mgr.db.Where("user_name = ? and service = ?", settings.UserName, settings.Service).Find(&existing).Error
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		return errors.New("notification settings doesn't exist")
	}

	const query = "update notification_settings set data = ? where user_name = ? and service = ?"
	return mgr.db.Exec(query, settings.Data, settings.UserName, settings.Service).Error
}
