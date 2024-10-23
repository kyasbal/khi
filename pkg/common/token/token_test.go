package token

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	token := New("foo")

	if token.RawToken != "foo" {
		t.Errorf("got %q, want %q", token.RawToken, "foo")
	}
	if !token.ValidAtLeastUntil.IsZero() {
		t.Errorf("got %q, want zero", token.ValidAtLeastUntil)
	}
}

func TestNewWithExpiry(t *testing.T) {
	expireTime := time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC)
	token := NewWithExpiry("foo", expireTime)

	if token.RawToken != "foo" {
		t.Errorf("got %q, want %q", token.RawToken, "foo")
	}
	if !expireTime.Equal(token.ValidAtLeastUntil) {
		t.Errorf("got %q, want %q", token.ValidAtLeastUntil, expireTime)
	}
}
