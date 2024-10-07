package structuredata

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func ReadReflect(sd StructureData, target interface{}) error {
	jsonStr, err := ToJson(sd)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), target)
}

func ReadReflectK8sManifest(sd StructureData, target runtime.Object) error {
	jsonStr, err := ToJson(sd)
	if err != nil {
		return err
	}
	scheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()
	_, _, err = deserializer.Decode([]byte(jsonStr), nil, target)
	if err != nil {
		return err
	}
	return nil
}
