package shared

import "errors"

var (
	ErrNoBooks   = errors.New("order must have at least one book")
	ErrNoCatalog = errors.New("book not found in catalog")
)
