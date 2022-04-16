package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken : creates a new token for the provided userID and duration
	CreateToken(userID uint, duration time.Duration) (string, *Payload, error)

	// VerifyToken : checks whether the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
