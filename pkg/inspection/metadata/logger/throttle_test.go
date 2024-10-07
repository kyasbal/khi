package logger

import (
	"testing"
)

func TestConstantLogThrottle(t *testing.T) {
	throttle := NewConstantLogThrottle(10)
	for i := 0; i < 9; i++ {
		if throttle.ThrottleStatus("foo") != StatusNoThrottle {
			t.Errorf("key foo shouldn't be throttled yet")
		}
		if throttle.ThrottleStatus("") != StatusNoThrottle {
			t.Errorf("key \"\" shouldn't be throttled never")
		}
	}
	if throttle.ThrottleStatus("foo") != StatusJustBeforeThrottle {
		t.Errorf("key foo should be throttled from next")
	}
	if throttle.ThrottleStatus("foo") != StatusThrottled {
		t.Errorf("foo should be throttled")
	}
	if throttle.ThrottleStatus("bar") != StatusNoThrottle {
		t.Errorf("key foo shouldn't be throttled yet")
	}
	if throttle.ThrottleStatus("") != StatusNoThrottle {
		t.Errorf("key \"\" shouldn't be throttled never")
	}
	if throttle.ThrottleStatus("") != StatusNoThrottle {
		t.Errorf("key \"\" shouldn't be throttled never")
	}
}
