package shrterr

import "errors"

var (
	ErrOriginalURLAlreadyExists = errors.New("original_url already exists")
)
