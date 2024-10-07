package enum

type LogType int

const (
	LogTypeUnknown               LogType = 0
	LogTypeEvent                 LogType = 1
	LogTypeAudit                 LogType = 2
	LogTypeContainer             LogType = 3
	LogTypeNode                  LogType = 4
	LogTypeGkeAudit              LogType = 5
	LogTypeComputeApi            LogType = 6
	LogTypeMulticloudAPI         LogType = 7
	LogTypeOnPremAPI             LogType = 8
	LogTypeNetworkAPI            LogType = 9
	LogTypeAutoscaler            LogType = 10
	LogTypeComposerEnvironment   LogType = 11
	LogTypeControlPlaneComponent LogType = 12

	logTypeUnusedEnd
)

type LogTypeFrontendMetadata struct {
	// EnumKeyName is the name of this enum value. Must match with the enum key.
	EnumKeyName string
	// Label string shown on frontnend to indicate the log type.
	Label string
	// Background color of the label on log pane.
	LabelBackgroundColor string
}

var LogTypes = map[LogType]LogTypeFrontendMetadata{
	LogTypeUnknown: {
		EnumKeyName:          "LogTypeUnknown",
		Label:                "unknown",
		LabelBackgroundColor: "#000000",
	},
	LogTypeEvent: {
		EnumKeyName:          "LogTypeEvent",
		Label:                "k8s_event",
		LabelBackgroundColor: "#3fb549",
	},
	LogTypeAudit: {
		EnumKeyName:          "LogTypeAudit",
		Label:                "k8s_audit",
		LabelBackgroundColor: "#000000",
	},
	LogTypeContainer: {
		EnumKeyName:          "LogTypeContainer",
		Label:                "k8s_container",
		LabelBackgroundColor: "#fe9bab",
	},
	LogTypeNode: {
		EnumKeyName:          "LogTypeNode",
		Label:                "k8s_node",
		LabelBackgroundColor: "#0077CC",
	},
	LogTypeGkeAudit: {
		EnumKeyName:          "LogTypeGkeAudit",
		Label:                "gke_audit",
		LabelBackgroundColor: "#AA00FF",
	},
	LogTypeComputeApi: {
		EnumKeyName:          "LogTypeComputeApi",
		Label:                "compute_api",
		LabelBackgroundColor: "#FFCC33",
	},
	LogTypeMulticloudAPI: {
		EnumKeyName:          "LogTypeMulticloudAPI",
		Label:                "multicloud_api",
		LabelBackgroundColor: "#AA00FF",
	},
	LogTypeOnPremAPI: {
		EnumKeyName:          "LogTypeOnPremAPI",
		Label:                "onprem_api",
		LabelBackgroundColor: "#AA00FF",
	},
	LogTypeNetworkAPI: {
		EnumKeyName:          "LogTypeNetworkAPI",
		Label:                "network_api",
		LabelBackgroundColor: "#33CCFF",
	},
	LogTypeAutoscaler: {
		EnumKeyName:          "LogTypeAutoscaler",
		Label:                "autoscaler",
		LabelBackgroundColor: "#FF5555",
	},
	LogTypeComposerEnvironment: {
		EnumKeyName:          "LogTypeComposerEnvironment",
		Label:                "composer_environment",
		LabelBackgroundColor: "#88AA55",
	},
	LogTypeControlPlaneComponent: {
		EnumKeyName:          "LogTypeControlPlaneComponent",
		Label:                "control_plane_component",
		LabelBackgroundColor: "#FF3333",
	},
}
