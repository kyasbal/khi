package model

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

var irregularPluralToSingularSuffixMap = map[string]string{
	"classes":    "class",
	"ingresses":  "ingress",
	"leases":     "lease",
	"dnses":      "dns",
	"identities": "identity",
	"policies":   "policy",
	"topologies": "topology",
}

type KubernetesObjectOperation struct {
	APIVersion      string
	PluralKind      string
	Namespace       string
	Name            string
	SubResourceName string
	Verb            enum.RevisionVerb
}

func (o *KubernetesObjectOperation) CovertToResourcePath() string {
	if o.SubResourceName != "" {
		return strings.ToLower(strings.Join([]string{
			o.APIVersion,
			o.GetSingularKindName(),
			o.Namespace,
			o.Name,
			o.SubResourceName,
		}, "#"))
	} else {
		return strings.ToLower(strings.Join([]string{
			o.APIVersion,
			o.GetSingularKindName(),
			o.Namespace,
			o.Name,
		}, "#"))
	}
}

func (o *KubernetesObjectOperation) GetSingularKindName() string {
	if strings.HasSuffix(o.PluralKind, "ses") || strings.HasSuffix(o.PluralKind, "ies") {
		for pluralSuffix, singularSuffix := range irregularPluralToSingularSuffixMap {
			if strings.HasSuffix(o.PluralKind, pluralSuffix) {
				return strings.TrimSuffix(o.PluralKind, pluralSuffix) + singularSuffix
			}
		}
		slog.Warn(fmt.Sprintf("unknown singular name for %s", o.PluralKind))
		return o.PluralKind
	}
	if strings.HasSuffix(o.PluralKind, "s") {
		return strings.TrimSuffix(o.PluralKind, "s")
	}
	slog.Error(fmt.Sprintf("unknown plural form %s!", o.PluralKind))
	return o.PluralKind
}
