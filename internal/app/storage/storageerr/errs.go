package storageerr

import "errors"

var (
	ErrEntityNotFound       = errors.New("entity not found")
	ErrNotUniqueShortcut    = errors.New("shortcut is not unique")
	ErrNotUniqueOriginalURL = errors.New("originalURL is not unique")
)
