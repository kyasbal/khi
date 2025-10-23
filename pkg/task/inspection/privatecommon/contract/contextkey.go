package privatecommon_contract

import (
	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
)

// APICallOptionsInjectorContextKey is the key to retrieve the list of googlecloud.CallOptionInjectorOption from context.
var APICallOptionsInjectorContextKey = typedmap.NewTypedKey[*[]googlecloud.CallOptionInjectorOption]("api-call-option-injector-options")
