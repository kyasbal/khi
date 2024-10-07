package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("plural:%s", tc.plural), func(t *testing.T) {
			o := KubernetesObjectOperation{PluralKind: tc.plural}
			assert.Equal(t, tc.singular, o.GetSingularKindName())
		})
	}
}
