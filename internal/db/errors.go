package db

import "errors"

var (
	ErrNotificationSettingsNotFound = errors.New("notification settings not found")

	ErrNotInTx     = errors.New("you're not in the tx")
	ErrAlreadyInTx = errors.New("you're already in the tx")
)
