package resourcepath

func ControlplaneComponent(cluster string, component string) string {
	if cluster == "" {
		cluster = nonSpecifiedPlaceholder
	}
	if component == "" {
		component = nonSpecifiedPlaceholder
	}
	return SubresourceLayerGeneralItem("@Cluster", "controlplane", "cluster-scope", cluster, component)
}
