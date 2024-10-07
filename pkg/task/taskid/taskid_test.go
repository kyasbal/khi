package taskid

import (
	"testing"
)

func TestNewTaskReference(t *testing.T) {
	referenceId := NewTaskReference("foo.bar")
	if referenceId.id != "foo.bar" {
		t.Errorf("Expected TaskReferenceId.id to be 'foo.bar', got %q", referenceId.id)
	}
}

func TestNewTaskImplementationId(t *testing.T) {
	testCases := []struct {
		taskId                 string
		wantId                 string
		wantImplementationHash string
	}{
		{taskId: "foo.bar", wantId: "foo.bar", wantImplementationHash: ""},
		{taskId: "hello#world", wantId: "hello", wantImplementationHash: "world"},
	}

	for _, tc := range testCases {
		implementationId := NewTaskImplementationId(tc.taskId)
		if implementationId.referenceId != tc.wantId {
			t.Errorf("NewTaskImplementationId(%q).id = %q, want %q", tc.taskId, implementationId.referenceId, tc.wantId)
		}
		if implementationId.implementationHash != tc.wantImplementationHash {
			t.Errorf("NewTaskImplementationId(%q).implementationHash = %q, want %q", tc.taskId, implementationId.implementationHash, tc.wantImplementationHash)
		}
	}
}

func TestTaskImplementationIdMatch(t *testing.T) {
	testCases := []struct {
		implementationId TaskImplementationId
		referenceId      TaskReferenceId
		want             bool
	}{
		{NewTaskImplementationId("foo"), NewTaskReference("foo"), true},
		{NewTaskImplementationId("hello#world"), NewTaskReference("hello"), true},
		{NewTaskImplementationId("abc"), NewTaskReference("xyz"), false},
	}

	for _, tc := range testCases {
		got := tc.implementationId.Match(tc.referenceId)
		if got != tc.want {
			t.Errorf("(%v).Match(%v) = %t, want %t", tc.implementationId, tc.referenceId, got, tc.want)
		}
	}
}

func TestDedupeReferenceIds(t *testing.T) {
	input := []TaskReferenceId{
		NewTaskReference("foo"),
		NewTaskReference("bar"),
		NewTaskReference("foo"), // Duplicate
		NewTaskReference("baz"),
	}

	deduped := DedupeReferenceIds(input)

	if len(deduped) != 3 {
		t.Fatalf("Expected 3 unique references, got %d", len(deduped))
	}

	expected := []string{"bar", "baz", "foo"} // Order after sorting
	for i, ref := range deduped {
		if ref.id != expected[i] {
			t.Errorf("Expected reference %d to have id %q, got %q", i, expected[i], ref.id)
		}
	}
}
