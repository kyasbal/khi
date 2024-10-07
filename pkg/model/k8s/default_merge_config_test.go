package k8s

import (
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
)

func TestGenerateDefaultMergeConfig(t *testing.T) {
	resolver, err := GenerateDefaultMergeConfig()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(resolver)
}

func TestBuilder(t *testing.T) {
	builder := appsv1.SchemeBuilder
	fmt.Println(builder)
}
