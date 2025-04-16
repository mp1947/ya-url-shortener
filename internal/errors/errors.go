package shrterr

import "errors"

var (
	ErrOriginalURLAlreadyExists     = errors.New("original_url already exists")
	ErrUnableToDetermineStorageType = errors.New("unable to determine storage type")
)
