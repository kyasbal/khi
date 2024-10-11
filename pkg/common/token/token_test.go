package token

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	token := New("foo")

	if token.RawToken != "foo" {
		t.Errorf("Expected token.RawToken to be 'foo', but got '%s'", token.RawToken)
	}
	if !token.ValidAtLeastUntil.IsZero() {
		t.Errorf("Expected token.ValidAtLeastUntil to be zero, but got '%s'", token.ValidAtLeastUntil)
	}
}

func TestNewWithExpiry(t *testing.T) {
	expireTime := time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC)
	token := NewWithExpiry("foo", expireTime)

	if token.RawToken != "foo" {
		t.Errorf("Expected token.RawToken to be 'foo', but got '%s'", token.RawToken)
	}
	if !expireTime.Equal(token.ValidAtLeastUntil) {
		t.Errorf("Expected token.ValidAtLeastUntil to be '%s', but got '%s'", expireTime, token.ValidAtLeastUntil)
	}
}
