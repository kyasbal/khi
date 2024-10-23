package model

import (
	"fmt"
	"testing"
)

func TestToSingularKindName(t *testing.T) {
	testCases := []struct {
		plural   string
		singular string
	}{
		{
			plural:   "pods",
			singular: "pod",
		},
		{
			plural:   "services",
			singular: "service",
		},
		{
			plural:   "ingresses",
			singular: "ingress",
		},
		{
			plural:   "clusterdnses",
			singular: "clusterdns",
		},
		{
			plural:   "csinodetopologies",
			singular: "csinodetopology",
		},
		{
			plural:   "entitlementidentities",
			singular: "entitlementidentity",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("plural:%s", tc.plural), func(t *testing.T) {
			o := KubernetesObjectOperation{PluralKind: tc.plural}

			if tc.singular != o.GetSingularKindName() {
				t.Errorf("got %q, want %q", o.GetSingularKindName(), tc.singular)
			}
		})
	}
}
