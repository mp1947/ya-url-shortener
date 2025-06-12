package shrterr

import "errors"

var (
	// ErrOriginalURLAlreadyExists is returned when an attempt is made to add a URL that already exists in the storage.
	ErrOriginalURLAlreadyExists = errors.New("original_url already exists")

	// ErrUnableToDetermineStorageType is returned when the application cannot identify or select a valid storage type for operation.
	ErrUnableToDetermineStorageType = errors.New("unable to determine storage type")
)
