package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	t.Run("returns error for empty email", func(t *testing.T) {
		require.ErrorIs(t, ValidateEmail(""), ErrEmailEmpty)
	})

	t.Run("accepts valid emails", func(t *testing.T) {
		emails := []string{
			"user@example.com",
			"first.last@example.co.uk",
			"user+tag@example.com",
			"user_name@sub-domain.example.org",
			"UPPER@EXAMPLE.COM",
			"a1@b2.io",
			"user-name@example-site.com",
		}

		for _, email := range emails {
			require.NoError(t, ValidateEmail(email), email)
		}
	})

	t.Run("rejects invalid emails", func(t *testing.T) {
		emails := []string{
			"plainaddress",
			"@example.com",
			"user@",
			"user@@example.com",
			"user@example",
			"user@.com",
			"user@example.",
			"user@-example.com",
			"user@example-.com",
			".user@example.com",
			"user.@example.com",
			"us..er@example.com",
			"user name@example.com",
			"user@exa mple.com",
			"user@example.com ",
			"юзер@пример.рф",
		}

		for _, email := range emails {
			require.ErrorIs(t, ValidateEmail(email), ErrEmailInvalid, email)
		}
	})

	t.Run("rejects email longer than max length", func(t *testing.T) {
		local := strings.Repeat("a", MaxEmailLength)
		require.ErrorIs(t, ValidateEmail(local+"@example.com"), ErrEmailInvalid)
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("returns error for empty password", func(t *testing.T) {
		require.ErrorIs(t, ValidatePassword(""), ErrPasswordEmpty)
	})

	t.Run("accepts valid passwords", func(t *testing.T) {
		passwords := []string{
			"passw0rd",
			"password1",
			"P@ssw0rd!",
			"super-secret-42",
			"Пароль123",
			"1234567a",
			strings.Repeat("a1", MaxPasswordLength/2),
		}

		for _, password := range passwords {
			require.NoError(t, ValidatePassword(password), password)
		}
	})

	t.Run("rejects too short password", func(t *testing.T) {
		require.ErrorIs(t, ValidatePassword("pass1"), ErrPasswordTooShort)
		require.ErrorIs(t, ValidatePassword(strings.Repeat("a1", MinPasswordLength/2-1)), ErrPasswordTooShort)
	})

	t.Run("rejects too long password", func(t *testing.T) {
		require.ErrorIs(t, ValidatePassword(strings.Repeat("a1", MaxPasswordLength/2+1)), ErrPasswordTooLong)
	})

	t.Run("rejects password without letter", func(t *testing.T) {
		require.ErrorIs(t, ValidatePassword("12345678"), ErrPasswordNoLetter)
	})

	t.Run("rejects password without digit", func(t *testing.T) {
		require.ErrorIs(t, ValidatePassword("password"), ErrPasswordNoDigit)
	})
}
