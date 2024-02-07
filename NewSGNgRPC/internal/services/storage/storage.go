package storage

import "errors"

var (
	ErrTextExist   = errors.New("Such text already exists")
	ErrIDNotFound  = errors.New("ID not found")
	ErrAppNotFound = errors.New("App not found")
)
