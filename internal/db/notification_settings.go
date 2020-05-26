package db

type NotificationSettings struct {
	ID       int    `json:"-"       db:"id"`
	UserName string `json:"-"       db:"user_name"`
	Service  string `json:"service" db:"service"`
	Data     string `json:"data"    db:"data"`
}

type NotificationSettingsMgr interface {
	FindNotificationSettings(userName string) ([]*NotificationSettings, error)
	FindNotificationSettingsForService(userName, service string) ([]*NotificationSettings, error)
	EnsureNotificationSettingsExists(settings *NotificationSettings) error
	UpdateNotificationSettings(settings *NotificationSettings) error
}

func (mgr *AppDatabaseMgr) FindNotificationSettings(userName string) ([]*NotificationSettings, error) {
	const query = "select * from notification_settings where user_name = $1"

	settings := []*NotificationSettings{}
	err := mgr.newdb.Select(&settings, query, userName)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (mgr *AppDatabaseMgr) FindNotificationSettingsForService(userName, service string) ([]*NotificationSettings, error) {
	const query = "select * from notification_settings where user_name = $1 and service = $2"

	settings := []*NotificationSettings{}
	err := mgr.newdb.Select(&settings, query, userName, service)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (mgr *AppDatabaseMgr) EnsureNotificationSettingsExists(settings *NotificationSettings) error {
	userSettings, err := mgr.FindNotificationSettingsForService(settings.UserName, settings.Service)
	if err != nil {
		return err
	}
	if len(userSettings) != 0 {
		return nil
	}

	const query = "insert into notification_settings (user_name, service, data) values ($1, $2, $3)"

	_, err = mgr.newdb.Exec(query, settings.UserName, settings.Service, settings.Data)

	return err
}

func (mgr *AppDatabaseMgr) UpdateNotificationSettings(settings *NotificationSettings) error {
	existing, err := mgr.FindNotificationSettingsForService(settings.UserName, settings.Service)
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		return ErrNotificationSettingsNotFound
	}

	const query = "update notification_settings set data = $1 where user_name = $2 and service = $3"

	_, err = mgr.newdb.Exec(query, settings.Data, settings.UserName, settings.Service)

	return err
}
