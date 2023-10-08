package storageerr

import "errors"

var (
	ErrEntityNotFound       = errors.New("Entity not found")
	ErrNotUniqueShortcut    = errors.New("shortcut is not unique")
	ErrNotUniqueOriginalURL = errors.New("OriginalURL is not unique")
)
