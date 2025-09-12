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

package googlecloudloggkeapiaudit_impl

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

func TestGkeAuditLogParser_ClusterCreationStartLog(t *testing.T) {
	userAccountName := "user@example.com"
	operationId := "operation-1726191072114-d3db4945-ad7b-4fff-aff7-55a867e4bc54"
	clusterResourcePath := resourcepath.Cluster("gke-basic-1")
	operationResourcePath := resourcepath.Operation(clusterResourcePath, "CreateCluster", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/cluster_creation_started.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantClusterRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbCreate,
		State:      enum.RevisionStateProvisioning,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2024-01-01T01:05:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, clusterResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, clusterResourcePath, wantClusterRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, clusterResourcePath, "cluster-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationStart,
		State:      enum.RevisionStateOperationStarted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2024-01-01T01:05:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_ClusterCreationFinishedLog(t *testing.T) {
	userAccountName := "user@example.com"
	operationId := "operation-1726191072114-d3db4945-ad7b-4fff-aff7-55a867e4bc54"
	clusterResourcePath := resourcepath.Cluster("gke-basic-1")
	operationResourcePath := resourcepath.Operation(clusterResourcePath, "CreateCluster", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/cluster_creation_started.yaml",
		&gkeAuditLogParser{}, nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantClusterRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbCreate,
		State:      enum.RevisionStateProvisioning,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2024-01-01T01:05:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, clusterResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, clusterResourcePath, wantClusterRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, clusterResourcePath, "cluster-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationStart,
		State:      enum.RevisionStateOperationStarted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2024-01-01T01:05:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_ClusterDeletionStartLog(t *testing.T) {
	userAccountName := "user@example.com"
	operationId := "operation-1726199159930-7409b104-8654-4667-b477-4ce504d09bea"
	clusterResourcePath := resourcepath.Cluster("gke-basic-1")
	operationResourcePath := resourcepath.Operation(clusterResourcePath, "DeleteCluster", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/cluster_deletion_started.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantClusterRevision := &history.StagingResourceRevision{
		State:      enum.RevisionStateDeleting,
		Verb:       enum.RevisionVerbDelete,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, clusterResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, clusterResourcePath, wantClusterRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, clusterResourcePath, "cluster-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationStart,
		State:      enum.RevisionStateOperationStarted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_ClusterDeletionFinishedLog(t *testing.T) {
	operationId := "operation-1726199159930-7409b104-8654-4667-b477-4ce504d09bea"
	userAccountName := "unknown"
	clusterResourcePath := resourcepath.Cluster("gke-basic-1")
	operationResourcePath := resourcepath.Operation(clusterResourcePath, "DeleteCluster", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/cluster_deletion_finished.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantClusterRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbDelete,
		State:      enum.RevisionStateDeleted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, clusterResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, clusterResourcePath, wantClusterRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, clusterResourcePath, "cluster-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationFinish,
		State:      enum.RevisionStateOperationFinished,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_NodepoolCreationStartLog(t *testing.T) {
	userAccountName := "user@example.com"
	operationId := "operation-1726191716103-f4072772-f902-453d-8776-b69047cebae6"
	nodepoolResourcePath := resourcepath.Nodepool("gke-basic-1", "default")
	operationResourcePath := resourcepath.Operation(nodepoolResourcePath, "CreateNodePool", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/nodepool_creation_started.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantNodepoolRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbCreate,
		State:      enum.RevisionStateProvisioning,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, nodepoolResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, nodepoolResourcePath, wantNodepoolRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, nodepoolResourcePath, "nodepool-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationStart,
		State:      enum.RevisionStateOperationStarted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_NodepoolCreationFinishedLog(t *testing.T) {
	userAccountName := "unknown"
	operationId := "operation-1726191716103-f4072772-f902-453d-8776-b69047cebae6"
	nodepoolResourcePath := resourcepath.Nodepool("gke-basic-1", "default")
	operationResourcePath := resourcepath.Operation(nodepoolResourcePath, "CreateNodePool", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/nodepool_creation_finished.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantNodepoolRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbCreate,
		State:      enum.RevisionStateExisting,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, nodepoolResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, nodepoolResourcePath, wantNodepoolRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, nodepoolResourcePath, "nodepool-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationFinish,
		State:      enum.RevisionStateOperationFinished,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_NodepoolDeletionStartLog(t *testing.T) {
	userAccountName := "user@example.com"
	operationId := "operation-1726191433631-f35aa16e-345f-4a0f-8091-ec613f0635c2"
	nodepoolResourcePath := resourcepath.Nodepool("gke-basic-1", "default-pool")
	operationResourcePath := resourcepath.Operation(nodepoolResourcePath, "DeleteNodePool", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/nodepool_deletion_started.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantNodepoolRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbDelete,
		State:      enum.RevisionStateDeleting,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, nodepoolResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, nodepoolResourcePath, wantNodepoolRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, nodepoolResourcePath, "nodepool-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationStart,
		State:      enum.RevisionStateOperationStarted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_NodepoolDeletionFinishedLog(t *testing.T) {
	userAccountName := "unknown"
	operationId := "operation-1726191433631-f35aa16e-345f-4a0f-8091-ec613f0635c2"
	nodepoolResourcePath := resourcepath.Nodepool("gke-basic-1", "default-pool")
	operationResourcePath := resourcepath.Operation(nodepoolResourcePath, "DeleteNodePool", operationId)

	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/nodepool_deletion_finished.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	wantNodepoolRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbDelete,
		State:      enum.RevisionStateDeleted,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, nodepoolResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, nodepoolResourcePath, wantNodepoolRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, nodepoolResourcePath, "nodepool-body")

	wantOperationRevision := &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbOperationFinish,
		State:      enum.RevisionStateOperationFinished,
		Requestor:  userAccountName,
		ChangeTime: testutil.MustParseTimeRFC3339("2025-01-01T00:00:00Z"),
	}
	history_test.AssertChangeSetHasCountOfRevisionsForResourcePath(t, cs, operationResourcePath, 1)
	history_test.AssertChangeSetHasRevisionForResourcePath(t, cs, operationResourcePath, wantOperationRevision, cmpopts.IgnoreFields(history.StagingResourceRevision{}, "Body"))
	history_test.AssertChangeSetHasRevisionMatchingBodyGoldensForResourcePath(t, cs, operationResourcePath, "operation-body")
}

func TestGkeAuditLogParser_ClusterCreationWithErrorLog(t *testing.T) {
	clusterName := "p0-gke-basic-1"
	clusterResourcePath := resourcepath.Cluster(clusterName)
	cs, err := parser_test.ParseFromYamlLogFile(
		"test/logs/gke_audit/cluster_creation_started_with_error.yaml",
		&gkeAuditLogParser{},
		nil, &gcpqueryutil.GCPCommonFieldSetReader{}, &gcpqueryutil.GCPMainMessageFieldSetReader{})
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	history_test.AssertChangeSetHasEventForResourcePath(t, cs, clusterResourcePath)
}
