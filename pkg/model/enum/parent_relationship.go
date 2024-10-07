package enum

type ParentRelationShip int

const (
	RelationshipChild                ParentRelationShip = 0
	RelationshipResourceStatus       ParentRelationShip = 1
	RelationshipOperation            ParentRelationShip = 2
	RelationshipEndpointSlice        ParentRelationShip = 3
	RelationshipContainer            ParentRelationShip = 4
	RelationshipNodeComponent        ParentRelationShip = 5
	RelationshipOwnerReference       ParentRelationShip = 6
	RelationshipPodBinding           ParentRelationShip = 7
	RelationshipNetworkEndpointGroup ParentRelationShip = 8
	RelationshipManagedInstanceGroup ParentRelationShip = 9
	relationshipUnusedEnd                               // Add items above. This field is used for counting items in this enum to test.
)

// parentRelationshipFrontendMetadata is a type defined for each parent relationship types.
type ParentRelationshipFrontendMetadata struct {
	Visible              bool
	EnumKeyName          string
	Label                string
	Hint                 string
	LabelColor           string
	LabelBackgroundColor string
}

var ParentRelationships = map[ParentRelationShip]ParentRelationshipFrontendMetadata{
	RelationshipChild: {
		Visible:              false,
		EnumKeyName:          "RelationshipChild",
		Label:                "",
		LabelColor:           "",
		LabelBackgroundColor: "",
	},
	RelationshipResourceStatus: {
		Visible:              true,
		EnumKeyName:          "RelationshipResourceStatus",
		Label:                "status",
		LabelColor:           "#FFFFFF",
		LabelBackgroundColor: "#4c29e8",
		Hint:                 "Resource status written on .status.conditions",
	},
	RelationshipOperation: {
		Visible:              true,
		EnumKeyName:          "RelationshipOperation",
		Label:                "operation",
		LabelColor:           "#FFFFFF",
		LabelBackgroundColor: "#000000",
		Hint:                 "GCP operations associated with this resource",
	},
	RelationshipEndpointSlice: {
		Visible:              true,
		EnumKeyName:          "RelationshipEndpointSlice",
		Label:                "endpointslice",
		LabelColor:           "#FFFFFF",
		LabelBackgroundColor: "#008000",
		Hint:                 "Pod serving status obtained from endpoint slice",
	},
	RelationshipContainer: {
		Visible:              true,
		EnumKeyName:          "RelationshipContainer",
		Label:                "container",
		LabelColor:           "#000000",
		LabelBackgroundColor: "#fe9bab",
		Hint:                 "Containers statuses/logs in Pods",
	},
	RelationshipNodeComponent: {
		Visible:              true,
		EnumKeyName:          "RelationshipNodeComponent",
		Label:                "node-component",
		LabelColor:           "#FFFFFF",
		LabelBackgroundColor: "#0077CC",
		Hint:                 "Non container resource running on a node",
	},
	RelationshipOwnerReference: {
		Visible:              true,
		EnumKeyName:          "RelationshipOwnerReference",
		Label:                "owns",
		LabelColor:           "#000000",
		LabelBackgroundColor: "#33DD88",
		Hint:                 "A k8s resource related to this resource from .metadata.ownerReference field",
	},
	RelationshipPodBinding: {
		Visible:              true,
		EnumKeyName:          "RelationshipPodBinding",
		Label:                "binds",
		LabelColor:           "#000000",
		LabelBackgroundColor: "#FF8855",
		Hint:                 "Pod binding subresource associated with this node",
	},
	RelationshipNetworkEndpointGroup: {
		Visible:              true,
		EnumKeyName:          "RelationshipNetworkEndpointGroup",
		Label:                "neg",
		LabelColor:           "#FFFFFF",
		LabelBackgroundColor: "#A52A2A",
		Hint:                 "Pod serving status obtained from the associated NEG status",
	},
	RelationshipManagedInstanceGroup: {
		Visible:              true,
		EnumKeyName:          "RelationshipManagedInstanceGroup",
		Label:                "mig",
		LabelColor:           "#FFFFFF",
		LabelBackgroundColor: "#FF5555",
		Hint:                 "MIG logs associated to the parent node pool",
	},
}
