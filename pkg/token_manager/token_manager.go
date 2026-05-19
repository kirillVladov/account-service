package token_manager

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired  = errors.New("token: expired")
	ErrTokenInvalid  = errors.New("token: invalid")
	ErrTokenRevoked  = errors.New("token: revoked")
	ErrTokenNotFound = errors.New("token: not found")
)

type Config struct {
	// Secret is the HMAC-SHA256 signing key for access tokens.
	Secret []byte

	AccessTTL  time.Duration
	RefreshTTL time.Duration

	// RefreshTokenBytes is the byte length of the random refresh token (default 32).
	RefreshTokenBytes int
}

type Manager struct {
	cfg Config
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func New(cfg Config) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

func (m *Manager) ValidateAccess(raw string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token: unexpected signing method: %v", t.Header["alg"])
		}

		return m.cfg.Secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrTokenInvalid
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (m *Manager) IssueAccess(userID string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.cfg.AccessTTL)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := t.SignedString(m.cfg.Secret)
	if err != nil {
		return "", fmt.Errorf("token: sign access: %w", err)
	}

	return signed, nil
}

func (m *Manager) IssueRefresh(userID string) (string, error) {
	buf := make([]byte, m.cfg.RefreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("token: generate refresh: %w", err)
	}

	raw := base64.URLEncoding.EncodeToString(buf)

	return raw, nil
}

func (m *Manager) IssuePair(userID, role string) (string, string, error) {
	access, err := m.IssueAccess(userID)
	if err != nil {
		return "", "", err
	}

	refresh, err := m.IssueRefresh(userID)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
