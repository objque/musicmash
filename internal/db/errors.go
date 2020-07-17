package db

import "errors"

var (
	ErrNotInTx     = errors.New("you're not in the tx")
	ErrAlreadyInTx = errors.New("you're already in the tx")
)
