package grouper

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	log_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/log"
)

func TestAllDependentLogGrouper(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*log.LogEntity
		wantKeys map[string]struct{}
	}{
		{
			name:     "empty logs",
			logs:     []*log.LogEntity{},
			wantKeys: map[string]struct{}{},
		},
		{
			name: "simple case",
			logs: []*log.LogEntity{
				log_test.MockLogWithId("id1"),
				log_test.MockLogWithId("id2"),
				log_test.MockLogWithId("id3"),
			},
			wantKeys: map[string]struct{}{
				"": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := AllDependentLogGrouper
			got := g.Group(tt.logs)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("Key length mismatch")
			}
			for wantKey := range tt.wantKeys {
				_, found := got[wantKey]
				if !found {
					t.Errorf("key %s was not found in the result", wantKey)
				}
			}
		})
	}
}

func TestAllIndependentLogGrouper(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*log.LogEntity
		wantKeys map[string]struct{}
	}{
		{
			name:     "empty logs",
			logs:     []*log.LogEntity{},
			wantKeys: map[string]struct{}{},
		},
		{
			name: "simple case",
			logs: []*log.LogEntity{
				log_test.MockLogWithId("id1"),
				log_test.MockLogWithId("id2"),
				log_test.MockLogWithId("id3"),
			},
			wantKeys: map[string]struct{}{
				"id1": {},
				"id2": {},
				"id3": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := AllIndependentLogGrouper
			got := g.Group(tt.logs)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("Key length mismatch")
			}
			for wantKey := range tt.wantKeys {
				_, found := got[wantKey]
				if !found {
					t.Errorf("key %s was not found in the result", wantKey)
				}
			}
		})
	}
}

func TestSingleStringFieldKeyLogGrouper(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*log.LogEntity
		wantKeys map[string]struct{}
	}{
		{
			name:     "empty logs",
			logs:     []*log.LogEntity{},
			wantKeys: map[string]struct{}{},
		},
		{
			name: "multiple logs",
			logs: []*log.LogEntity{
				log_test.MustLogEntity("textPayload: log message 1\nkey: groupA"),
				log_test.MustLogEntity("textPayload: log message 2\nkey: groupB"),
				log_test.MustLogEntity("textPayload: log message 3\nkey: groupA"),
			},
			wantKeys: map[string]struct{}{
				"groupA": {},
				"groupB": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewSingleStringFieldKeyLogGrouper("key")
			got := g.Group(tt.logs)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("Key length mismatch")
			}
			for wantKey := range tt.wantKeys {
				_, found := got[wantKey]
				if !found {
					t.Errorf("key %s was not found in the result", wantKey)
				}
			}
		})
	}
}
