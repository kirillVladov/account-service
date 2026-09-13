package errs

import "errors"

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountBlocked       = errors.New("account blocked")
	ErrTokenNotValid        = errors.New("token not valid")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)
