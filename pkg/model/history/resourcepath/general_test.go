package resourcepath

import (
	"testing"
)

func TestAPIVersionLayerGeneralItem(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		want       string
	}{
		{
			name:       "basic",
			apiVersion: "core/v1",
			want:       "core/v1",
		},
		{
			name:       "core package can omit the package domain",
			apiVersion: "v1",
			want:       "core/v1",
		},
		{
			name:       "empty",
			apiVersion: "",
			want:       "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := APIVersionLayerGeneralItem(tt.apiVersion); got != tt.want {
				t.Errorf("APIVersionLayerGeneralItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKindLayerGeneralItem(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		kind       string
		want       string
	}{
		{
			name:       "basic",
			apiVersion: "v1",
			kind:       "Pod",
			want:       "core/v1#Pod",
		},
		{
			name:       "empty apiVersion",
			apiVersion: "",
			kind:       "Pod",
			want:       "unknown#Pod",
		},
		{
			name:       "empty kind",
			apiVersion: "v1",
			kind:       "",
			want:       "core/v1#unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KindLayerGeneralItem(tt.apiVersion, tt.kind); got != tt.want {
				t.Errorf("KindLayerGeneralItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceLayerGeneralItem(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		kind       string
		namespace  string
		want       string
	}{
		{
			name:       "basic",
			apiVersion: "v1",
			kind:       "Pod",
			namespace:  "default",
			want:       "core/v1#Pod#default",
		},
		{
			name:       "empty apiVersion",
			apiVersion: "",
			kind:       "Pod",
			namespace:  "default",
			want:       "unknown#Pod#default",
		},
		{
			name:       "empty kind",
			apiVersion: "v1",
			kind:       "",
			namespace:  "default",
			want:       "core/v1#unknown#default",
		},
		{
			name:       "empty namespace",
			apiVersion: "v1",
			kind:       "Pod",
			namespace:  "",
			want:       "core/v1#Pod#unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NamespaceLayerGeneralItem(tt.apiVersion, tt.kind, tt.namespace); got != tt.want {
				t.Errorf("NamespaceLayerGeneralItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNameLayerGeneralItem(t *testing.T) {
	tests := []struct {
		tname      string
		apiVersion string
		kind       string
		namespace  string
		name       string
		want       string
	}{
		{
			tname:      "basic",
			apiVersion: "v1",
			kind:       "Pod",
			namespace:  "default",
			name:       "my-pod",
			want:       "core/v1#Pod#default#my-pod",
		},
		{
			tname:      "empty apiVersion",
			apiVersion: "",
			kind:       "Pod",
			namespace:  "default",
			name:       "my-pod",
			want:       "unknown#Pod#default#my-pod",
		},
		{
			tname:      "empty kind",
			apiVersion: "v1",
			kind:       "",
			namespace:  "default",
			name:       "my-pod",
			want:       "core/v1#unknown#default#my-pod",
		},
		{
			tname:      "empty namespace",
			apiVersion: "v1",
			kind:       "Pod",
			namespace:  "",
			name:       "my-pod",
			want:       "core/v1#Pod#unknown#my-pod",
		},
		{
			tname:      "empty name",
			apiVersion: "v1",
			kind:       "Pod",
			namespace:  "default",
			name:       "",
			want:       "core/v1#Pod#default#unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.tname, func(t *testing.T) {
			if got := NameLayerGeneralItem(tt.apiVersion, tt.kind, tt.namespace, tt.name); got != tt.want {
				t.Errorf("NameLayerGeneralItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubresourceLayerGeneralItem(t *testing.T) {
	tests := []struct {
		tname       string
		apiVersion  string
		kind        string
		namespace   string
		name        string
		subresource string
		want        string
	}{
		{
			tname:       "basic",
			apiVersion:  "v1",
			kind:        "Pod",
			namespace:   "default",
			name:        "my-pod",
			subresource: "status",
			want:        "core/v1#Pod#default#my-pod#status",
		},
		{
			tname:       "empty apiVersion",
			apiVersion:  "",
			kind:        "Pod",
			namespace:   "default",
			name:        "my-pod",
			subresource: "status",
			want:        "unknown#Pod#default#my-pod#status",
		},
		{
			tname:       "empty kind",
			apiVersion:  "v1",
			kind:        "",
			namespace:   "default",
			name:        "my-pod",
			subresource: "status",
			want:        "core/v1#unknown#default#my-pod#status",
		},
		{
			tname:       "empty namespace",
			apiVersion:  "v1",
			kind:        "Pod",
			namespace:   "",
			name:        "my-pod",
			subresource: "status",
			want:        "core/v1#Pod#unknown#my-pod#status",
		},
		{
			tname:       "empty name",
			apiVersion:  "v1",
			kind:        "Pod",
			namespace:   "default",
			name:        "",
			subresource: "status",
			want:        "core/v1#Pod#default#unknown#status",
		},
		{
			tname:       "empty subresource",
			apiVersion:  "v1",
			kind:        "Pod",
			namespace:   "default",
			name:        "my-pod",
			subresource: "",
			want:        "core/v1#Pod#default#my-pod#unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.tname, func(t *testing.T) {
			if got := SubresourceLayerGeneralItem(tt.apiVersion, tt.kind, tt.namespace, tt.name, tt.subresource); got != tt.want {
				t.Errorf("SubresourceLayerGeneralItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
