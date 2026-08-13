package errs

import "errors"

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrForbidden            = errors.New("forbidden")
)
