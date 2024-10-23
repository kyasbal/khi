package resourcepath

// NetworkEndpointGroup returns the ResourcePath of timeline for NEG.
func NetworkEndpointGroup(negNamespace string, negName string) ResourcePath {
	if negNamespace == "" {
		negNamespace = nonSpecifiedPlaceholder
	}
	if negName == "" {
		negName = nonSpecifiedPlaceholder
	}
	return NameLayerGeneralItem("networking.gke.io/v1beta1", "servicenetworkendpointgroup", negNamespace, negName)
}
