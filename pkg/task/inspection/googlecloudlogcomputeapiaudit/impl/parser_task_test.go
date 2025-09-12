// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package googlecloudlogcomputeapiaudit_impl

import (
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/gcpqueryutil"
	"github.com/GoogleCloudPlatform/khi/pkg/model/enum"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history/resourcepath"
	history_test "github.com/GoogleCloudPlatform/khi/pkg/model/history/test"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil"
	parser_test "github.com/GoogleCloudPlatform/khi/pkg/testutil/parser"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestComputeApiParser_Parse_OperationFirstLog(t *testing.T) {
	nodeName := "gke-gke-basic-1-default-5e5b794d-2m33"
	serviceAccountName := "serviceaccount@project-id.iam.gserviceaccount.com"
	operationId := "operation-1726191739294-621f6556f5492-0777bde4-78d02b5a"
	wantLogSummary := "v1.compute.instances.insert Started"
	nodeResourcePath := resourcepath.Node(nodeName)
	operationResourcePath := resourcepath.Operation(nodeResourcePath, "insert", operationId)
	wantRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationStart,
		State:      enum.RevisionStateOperationStarted,
		Requestor:  serviceAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2024-01-01T01:00:00Z"),
		Partial:    false,
	}

	cs, err := parser_test.ParseFromYamlLogFile("test/logs/compute_api/operation_first.yaml", &computeAPIParser{}, nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	history_test.AssertChangeSetHasLogSummary(t, cs, wantLogSummary)

	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "request")

	history_test.AssertChangeSetHasEventForResourcePath(t, cs, nodeResourcePath)
}

func TestComputeApiParser_Parse_OperationLastLog(t *testing.T) {
	nodeName := "gke-gke-basic-1-default-5e5b794d-2m33"
	serviceAccountName := "serviceaccount@project-id.iam.gserviceaccount.com"
	operationId := "operation-1726191739294-621f6556f5492-0777bde4-78d02b5a"
	wantLogSummary := "v1.compute.instances.insert Finished"
	nodeResourcePath := resourcepath.Node(nodeName)
	operationResourcePath := resourcepath.Operation(nodeResourcePath, "insert", operationId)
	wantRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationFinish,
		State:      enum.RevisionStateOperationFinished,
		Requestor:  serviceAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2024-01-01T01:05:00Z"),
		Partial:    false,
	}

	cs, err := parser_test.ParseFromYamlLogFile("test/logs/compute_api/operation_last.yaml", &computeAPIParser{}, nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	history_test.AssertChangeSetHasLogSummary(t, cs, wantLogSummary)

	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "request")

	history_test.AssertChangeSetHasEventForResourcePath(t, cs, nodeResourcePath)
}

func TestComputeApiParser_ParseOperationFirstLastLog(t *testing.T) {
	nodeName := "gke-basic-test-abcd"
	wantLogSummary := "compute.instances.repair.recreateInstance"
	nodeResourcePath := resourcepath.Node(nodeName)

	cs, err := parser_test.ParseFromYamlLogFile("test/logs/compute_api/operation_firstlast.yaml", &computeAPIParser{}, nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	history_test.AssertChangeSetHasLogSummary(t, cs, wantLogSummary)
	history_test.AssertChangeSetHasEventForResourcePath(t, cs, nodeResourcePath)
}
