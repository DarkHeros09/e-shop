package token

import "time"

// Maker is an interface for manging tokens
type Maker interface {
	// CreateToken creates a new token for specific username and duration
	CreateToken(userID int64, username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is vslid or not
	VerifyToken(token string) (*Payload, error)
}
