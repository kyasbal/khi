package model

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

var irregularPluralToSingularEndWithSes = map[string]string{
	"ingresses": "ingress",
	"leases":    "lease",
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
	if strings.HasSuffix(o.PluralKind, "ses") { // leases -> lease, ingresses -> ingress
		if strings.HasSuffix(o.PluralKind, "classes") { // for priorityclasses,storageclasses,runtimeclasses
			return strings.TrimSuffix(o.PluralKind, "es")
		}
		if singular, found := irregularPluralToSingularEndWithSes[o.PluralKind]; found {
			return singular
		} else {
			slog.Warn(fmt.Sprintf("unknown singular name for %s", o.PluralKind))
			return o.PluralKind
		}
	}
	if strings.HasSuffix(o.PluralKind, "s") {
		return strings.TrimSuffix(o.PluralKind, "s")
	}
	slog.Error(fmt.Sprintf("unknown plural form %s!", o.PluralKind))
	return o.PluralKind
}
