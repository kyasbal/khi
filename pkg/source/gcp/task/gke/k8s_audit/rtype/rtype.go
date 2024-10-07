package rtype

// rtype.Type indicates the schema of log bodies under request or response.
type Type = int

const (
	RTypeUnknown       Type = 0
	RTypePatch         Type = 1
	RTypeDeleteOptions Type = 2
	RTypeSatus         Type = 3
	RTypeUnusedEnd
)

var Types map[string]Type = map[string]Type{
	"k8s.io/Patch":                         RTypePatch,
	"meta.k8s.io/__internal.DeleteOptions": RTypeDeleteOptions,
	"core.k8s.io/v1.Status":                RTypeSatus,
}
