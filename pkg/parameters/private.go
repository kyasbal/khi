package parameters

import (
	"flag"
	"strings"
)

var Private *PrivateParameters = &PrivateParameters{}

type PrivateParameters struct {
	// InspectionMode
	// If this flag is set, KHI use the inspection token to access APIs.
	InspectionMode *bool

	// IAMToken is the initial IAM token used for accessing APIs.
	// InspectionMode will be true when only this value is specified.
	IAMToken *string

	// GALabels is the additional labels sent with the other labels. This flag is needed for taking analytics.
	// The expected format is "key1=value1,key2=value2,..."
	GALabels *string
}

// PostProcess implements ParameterStore.
func (p *PrivateParameters) PostProcess() error {
	if p.IAMToken != nil && *p.IAMToken != "" {
		*p.InspectionMode = true
	}
	return nil
}

// Prepare implements ParameterStore.
func (p *PrivateParameters) Prepare() error {
	p.InspectionMode = flag.Bool("inspection-mode", false, "If this flag is set, KHI use the inspection token to access APIs.")
	p.IAMToken = flag.String("iam-token", "", "The initial IAM token used for accessing APIs. InspectionMode will be true when only this value is specified.")
	p.GALabels = flag.String("ga-labels", "", "The additional labels sent with the other labels. This flag is needed for taking analytics.")
	return nil
}

// GetMapOfGALabels convert the GALabels parameter into a map.
// The expected format is "key1=value1,key2=value2,..."
func (p *PrivateParameters) GetMapOfGALabels() map[string]string {
	result := make(map[string]string)
	keyValuePairs := strings.Split(*p.GALabels, ",")
	for _, pair := range keyValuePairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		key, value, found := strings.Cut(pair, "=")
		if !found {
			value = "null"
		}
		result[key] = value
	}
	return result
}

var _ ParameterStore = (*PrivateParameters)(nil)
