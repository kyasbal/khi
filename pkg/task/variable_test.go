package task

import (
	"context"
	"testing"
	"time"
)

func TestVariableSetToBeReleasedWhenItWasSet(t *testing.T) {
	variableSet := NewVariableSet(map[string]any{})
	go func() {
		<-time.After(time.Second)
		variableSet.Set("foo-key", "bar")
	}()
	val, err := variableSet.Wait(context.Background(), "foo-key")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if val != "bar" {
		t.Errorf("variable `foo-key` wasn't resolved to `bar`")
	}
}

func TestVariableSetToBeCancelledWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	variableSet := NewVariableSet(map[string]any{})
	go func() {
		<-time.After(time.Second)
		cancel()
	}()
	_, err := variableSet.Wait(ctx, "foo-key")
	if err != context.Canceled {
		t.Errorf("unexpected error %v", err)
	}
	_, err = variableSet.Wait(ctx, "foo-key")
	if err != context.Canceled {
		t.Errorf("unexpected error %v", err)
	}
}

func TestVariableSetShouldReturnErrorWhenTheVariableSetTwice(t *testing.T) {
	variableSet := NewVariableSet(map[string]any{})
	err1 := variableSet.Set("foo-key", "bar1")
	err2 := variableSet.Set("foo-key", "bar2")
	if err1 != nil {
		t.Errorf("err1 is expected to be nil. but an error occured.")
	}
	if err2 == nil {
		t.Errorf("err2 is expected to be an error. but no error occured.")
	}
}

func TestGetTypedVariableFromTaskVariable(t *testing.T) {
	vs := NewVariableSet(map[string]any{})
	err := vs.Set("foo", time.Date(2023, time.April, 1, 1, 1, 1, 1, time.UTC))
	if err != nil {
		t.Errorf("unexpected error\n%s", err)
	}
	result, err := GetTypedVariableFromTaskVariable[time.Time](vs, "foo", time.Time{})
	if err != nil {
		t.Errorf("unexpected error\n%s", err)
	}
	if result.String() != time.Date(2023, time.April, 1, 1, 1, 1, 1, time.UTC).String() {
		t.Errorf("not matching with the expected value\n%s", err)
	}
}

func TestGetTypedVariableFromTaskVariableWithInvalidType(t *testing.T) {
	vs := NewVariableSet(map[string]any{})
	err := vs.Set("foo", time.Date(2023, time.April, 1, 1, 1, 1, 1, time.UTC))
	if err != nil {
		t.Errorf("unexpected error\n%s", err)
	}
	_, err = GetTypedVariableFromTaskVariable[*time.Time](vs, "foo", nil)
	if err.Error() != "the given value 2023-04-01 01:01:01.000000001 +0000 UTC in foo couldn't be converted to *time.Time" {
		t.Errorf("the result error is not matching with the expected error\n%s", err)
	}
}
