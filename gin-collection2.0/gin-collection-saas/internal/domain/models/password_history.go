package models

import "time"

// PasswordHistory stores previous password hashes for a user
// Used to prevent password reuse
type PasswordHistory struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	CreatedAt    time.Time `json:"created_at"`
}

// MaxPasswordHistory is the number of previous passwords to remember
const MaxPasswordHistory = 5
