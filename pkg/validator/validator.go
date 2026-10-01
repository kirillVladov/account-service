package validator

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
	MaxEmailLength    = 254
)

var (
	ErrEmailEmpty       = errors.New("email is empty")
	ErrEmailInvalid     = errors.New("email is invalid")
	ErrPasswordEmpty    = errors.New("password is empty")
	ErrPasswordTooShort = errors.New("password is too short")
	ErrPasswordTooLong  = errors.New("password is too long")
	ErrPasswordNoLetter = errors.New("password must contain at least one letter")
	ErrPasswordNoDigit  = errors.New("password must contain at least one digit")
)

var emailRe = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@` +
		`[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?` +
		`(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`,
)

func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailEmpty
	}

	if utf8.RuneCountInString(email) > MaxEmailLength {
		return ErrEmailInvalid
	}

	local, _, found := strings.Cut(email, "@")
	if !found || strings.HasPrefix(local, ".") || strings.HasSuffix(local, ".") || strings.Contains(local, "..") {
		return ErrEmailInvalid
	}

	if !emailRe.MatchString(email) {
		return ErrEmailInvalid
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	length := utf8.RuneCountInString(password)
	if length < MinPasswordLength {
		return ErrPasswordTooShort
	}

	if length > MaxPasswordLength {
		return ErrPasswordTooLong
	}

	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasLetter {
		return ErrPasswordNoLetter
	}

	if !hasDigit {
		return ErrPasswordNoDigit
	}

	return nil
}
