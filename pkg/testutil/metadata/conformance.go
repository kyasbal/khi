package metadata_test

import (
	"encoding/json"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
)

func ConformanceMetadataTypeTest(t *testing.T, m metadata.Metadata) {
	t.Run("metadata type must be serializable", func(t *testing.T) {
		ConformanceTestMetadataIsSerializable(t, m)
	})
}

func ConformanceTestMetadataIsSerializable(t *testing.T, m metadata.Metadata) {
	_, err := json.Marshal(m.ToSerializable())
	if err != nil {
		t.Errorf("Expected metadata is JSON serializable. But returned an error\n%v", err)
	}
}
