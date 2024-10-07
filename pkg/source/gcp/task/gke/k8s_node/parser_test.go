package k8s_node

import (
	"testing"

	log_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/log"
	"github.com/google/go-cmp/cmp"
)

func TestParseSummary(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{{
		name: "parse pod status",
		input: `jsonPayload:
  MESSAGE: I0101 22:34:24.745733    1754 kubelet_getters.go:187] "Pod status updated" pod="kube-system/kube-proxy-foo" status=Running`,
		expected: "Pod status updated(status=Running)",
	}, {
		name: "PLEG event",
		input: `jsonPayload:
  MESSAGE: 'I0101 09:35:31.624702    1754 kubelet.go:2375] "SyncLoop (PLEG): event for pod" pod="kube-system/fluentbit-gke-rxrcn" event=&{ID:caecee7e-19ba-4463-aa9d-2a46275c077c Type:ContainerStarted Data:4a4f3a5cc34c7d2f39bdc22690b8a96b1241d4367c5c147f09413c44101755e9}'`,
		expected: "SyncLoop (PLEG): event for pod(ContainerStarted)",
	}, {
		name: "node event",
		input: `jsonPayload:
  MESSAGE: I0101 09:31:25.067393    1761 kubelet_node_status.go:669] "Recording event message for node" node="node name" event="NodeHasSufficientMemory"`,
		expected: "Recording event message for node(NodeHasSufficientMemory)",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := log_test.MustLogEntity(tc.input)
			summary, err := parseSummary(l)
			if err != nil {
				t.Fatal(err)
			}
			if summary != tc.expected {
				t.Errorf("expected %s,\nactual %s", tc.expected, summary)
			}
		})
	}
}

func TestParseRunPodSandboxLog(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected *runPodSandboxLog
	}{
		{
			Name:  "standard run pod log",
			Input: "RunPodSandbox for &PodSandboxMetadata{Name:podname,Uid:b86b49f2431d244c613996c6472eb864,Namespace:kube-system,Attempt:0,} returns sandbox id \"6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1\"",
			Expected: &runPodSandboxLog{
				PodName:      "podname",
				PodNamespace: "kube-system",
				PodSandboxId: "6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1",
			},
		},
		{
			Name:     "missing container id",
			Input:    "RunPodSandbox for &PodSandboxMetadata{Name:podname,Uid:b86b49f2431d244c613996c6472eb864,Namespace:kube-system,Attempt:0,}",
			Expected: &runPodSandboxLog{PodName: "podname", PodNamespace: "kube-system"},
		},
		{
			Name:     "missing pod name in the struct",
			Input:    "RunPodSandbox for &PodSandboxMetadata{Uid:b86b49f2431d244c613996c6472eb864,Namespace:kube-system,Attempt:0,} returns sandbox id \"6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1\"",
			Expected: nil,
		},
		{
			Name:     "missing pod metadata struct",
			Input:    "RunPodSandbox for returns sandbox id \"6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1\"",
			Expected: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			sandbox, _ := parseRunPodSandboxLog(tc.Input)
			if diff := cmp.Diff(tc.Expected, sandbox); diff != "" {
				t.Errorf("the result sandbox is not matching with the expected result\n%s", diff)
			}
		})
	}
}

func TestParseCreateContainerLog(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected *createContainerLog
	}{
		{
			Name:  "standard create container log",
			Input: "CreateContainer within sandbox \"e175052cada9b999c5d9fabc8dc2276effc92b564aff74633eee122bcd4c8097\" for &ContainerMetadata{Name:config-init,Attempt:0,} returns container id \"14a996c61131027c75cc9e454acd8244c23ff7ddd236ee4ebbd0dd18d7d637d8\"",
			Expected: &createContainerLog{
				PodSandboxId:  "e175052cada9b999c5d9fabc8dc2276effc92b564aff74633eee122bcd4c8097",
				ContainerId:   "14a996c61131027c75cc9e454acd8244c23ff7ddd236ee4ebbd0dd18d7d637d8",
				ContainerName: "config-init",
			},
		},
		{
			Name:  "standard create container log without container id",
			Input: "CreateContainer within sandbox \"e175052cada9b999c5d9fabc8dc2276effc92b564aff74633eee122bcd4c8097\" for &ContainerMetadata{Name:config-init,Attempt:0,}",
			Expected: &createContainerLog{
				PodSandboxId:  "e175052cada9b999c5d9fabc8dc2276effc92b564aff74633eee122bcd4c8097",
				ContainerId:   "",
				ContainerName: "config-init",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			createContainerLog, _ := parseCreateContainerLog(tc.Input)
			if diff := cmp.Diff(tc.Expected, createContainerLog); diff != "" {
				t.Errorf("the result  is not matching with the expected result\n%s", diff)
			}
		})
	}
}

func TestReadGoStructFromString(t *testing.T) {
	testCases := []struct {
		Name       string
		Input      string
		StructName string
		Expected   map[string]string
	}{
		{
			Name:       "An example RunPodSandbox log",
			Input:      "RunPodSandbox for &PodSandboxMetadata{Name:podname,Uid:b86b49f2431d244c613996c6472eb864,Namespace:kube-system,Attempt:0,} returns sandbox id \"6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1\"",
			StructName: "PodSandboxMetadata",
			Expected: map[string]string{
				"Name":      "podname",
				"Namespace": "kube-system",
				"Attempt":   "0",
				"Uid":       "b86b49f2431d244c613996c6472eb864",
			},
		},
		{
			Name:       "An example CreateContainer log",
			Input:      "CreateContainer within sandbox \"573208ed2827243aa3db0db52e8f5a8d6fe65fcf67d93ecc76f5a4d92378af83\" for &ContainerMetadata{Name:fluentbit-gke-init,Attempt:0,} returns container id \"fc3e6702e38e918ec02567358c4c889b38fc628838645222d9a08b0b68c90256\"",
			StructName: "ContainerMetadata",
			Expected: map[string]string{
				"Attempt": "0",
				"Name":    "fluentbit-gke-init",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := readGoStructFromString(tc.Input, tc.StructName)
			if diff := cmp.Diff(tc.Expected, result); diff != "" {
				t.Errorf("result is not matching with the expected result\n%s", diff)
			}
		})
	}
}

func TestSafeParseContainerId(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "standard container id",
			Input:    "4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240",
			Expected: "4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240",
		},
		{
			Name:     "with containerd scheme",
			Input:    "containerd://4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240",
			Expected: "4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240",
		},
		{
			Name:     "json like text",
			Input:    "{Type:containerd ID:4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240}",
			Expected: "4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240",
		},
		{
			Name:     "json like text with double quotes",
			Input:    "{\"Type\":\"containerd\",\"ID\":\"26f7b5144c6b7edfd132852e6bfdcdcb5cbd5881cfb0512f8dd1b5b3d33c31de\"}",
			Expected: "26f7b5144c6b7edfd132852e6bfdcdcb5cbd5881cfb0512f8dd1b5b3d33c31de",
		},
		{
			Name:     "json like text changed order",
			Input:    "{ID:4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240 Type:containerd}",
			Expected: "4f376313b00249b84d4783cefc3e88fa59ac0362209937bcc2cdedfb97724240",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := safeParseContainerId(tc.Input)
			if diff := cmp.Diff(tc.Expected, actual); diff != "" {
				t.Errorf("the parsed container id is not matching with the expected value\n%s", diff)
			}
		})
	}
}

func TestReadNextQuotedString(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "standard input obtained from RunPodSandbox",
			Input:    "returns sandbox id \"6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1\"",
			Expected: "6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1",
		},
		{
			Name:     "not containing quote",
			Input:    "foo bar",
			Expected: "",
		},
		{
			Name:     "contains single double quote",
			Input:    "\"foo bar",
			Expected: "",
		},
		{
			Name:     "contains more than 3 double quote",
			Input:    "\"foo bar\" \"qux baz\"",
			Expected: "foo bar",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			nextQuoted := readNextQuotedString(tc.Input)
			if nextQuoted != tc.Expected {
				t.Errorf("expected:%s\nactual:%s", tc.Expected, nextQuoted)
			}
		})
	}
}

func TestRewroteContainerId(t *testing.T) {
	testCases := []struct {
		Name            string
		Id              string
		ReadableName    string
		OriginalMessage string
		ExpectedMessage string
	}{
		{
			Name:            "with full length of Id",
			OriginalMessage: "CreateContainer within sandbox \"8ff41ed8b310695c2223a702261c94d33675e1ac03442b2dc73b06ed11478f32\" for container &ContainerMetadata{Name:repeat-ready,Attempt:0,}",
			ExpectedMessage: "CreateContainer within sandbox \"8ff41ed...(1-1-probes/repeat-ready)\" for container &ContainerMetadata{Name:repeat-ready,Attempt:0,}",
			ReadableName:    "1-1-probes/repeat-ready",
			Id:              "8ff41ed8b310695c2223a702261c94d33675e1ac03442b2dc73b06ed11478f32",
		},
		{
			Name:            "without Id",
			OriginalMessage: "CreateContainer within sandbox \"8ff41ed8b310695c2223a702261c94d33675e1ac03442b2dc73b06ed11478f32\" for container &ContainerMetadata{Name:repeat-ready,Attempt:0,}",
			ExpectedMessage: "CreateContainer within sandbox \"8ff41ed8b310695c2223a702261c94d33675e1ac03442b2dc73b06ed11478f32\" for container &ContainerMetadata{Name:repeat-ready,Attempt:0,}",
			ReadableName:    "",
			Id:              "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			msg := rewriteIdWithReadableName(tc.Id, tc.ReadableName, tc.OriginalMessage)
			if msg != tc.ExpectedMessage {
				t.Errorf("expected:%s\nactual:%s", tc.ExpectedMessage, msg)
			}
		})
	}
}
