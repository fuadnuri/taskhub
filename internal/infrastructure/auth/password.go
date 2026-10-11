package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const defaultCost = bcrypt.DefaultCost

// PasswordHasher wraps bcrypt hashing logic.
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a PasswordHasher with the given cost.
// Pass 0 to use bcrypt.DefaultCost.
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost == 0 {
		cost = defaultCost
	}
	return &PasswordHasher{cost: cost}
}

// Hash returns a bcrypt hash of the plain-text password.
func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(bytes), nil
}

// Verify returns nil if password matches hash, or an error otherwise.
func (h *PasswordHasher) Verify(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	return nil
}
