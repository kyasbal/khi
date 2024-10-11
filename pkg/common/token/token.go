package token

import "time"

// Token contains the token string like access token.
type Token struct {
	// RawToken is the actual token in string representation.
	RawToken string
	// ValidAtLeastUntil holds the expiration time when the information is available.
	// The token refreshers won't refresh token if this value exists and later than now.
	// The default value will be `January 1, year 1, 00:00:00.000000000 UTC` and it is earlier than the possible time.Now().
	// The default value naturally means the token is expired.
	ValidAtLeastUntil time.Time
}

// NewWithExpiry instanciate a new Token from the raw token string and expiration time.
func NewWithExpiry(rawToken string, validAtLeastUntil time.Time) *Token {
	return &Token{
		RawToken:          rawToken,
		ValidAtLeastUntil: validAtLeastUntil,
	}
}

// New instanciate a new Token without expiration time.
func New(rawToken string) *Token {
	return &Token{
		RawToken: rawToken,
	}
}
