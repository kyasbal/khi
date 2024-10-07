package k8s_audit

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/bindingrecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/commonrecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/containerstatusrecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/endpointslicerecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/noderecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/ownerreferencerecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/snegrecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder/statusrecorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/v2commonlogparse"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/v2logconvert"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/v2manifestgenerate"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/v2timelinegrouping"
)

var PrepareK8sAuditTasks inspection.PrepareInspectionServerFunc = func(inspectionServer *inspection.InspectionTaskServer) error {
	err := inspectionServer.AddTaskDefinition(v2commonlogparse.Task)
	if err != nil {
		return err
	}

	err = inspectionServer.AddTaskDefinition(v2timelinegrouping.Task)
	if err != nil {
		return err
	}

	err = inspectionServer.AddTaskDefinition(v2manifestgenerate.Task)
	if err != nil {
		return err
	}

	err = inspectionServer.AddTaskDefinition(v2logconvert.Task)
	if err != nil {
		return err
	}
	manager := recorder.NewTaskManager()
	err = commonrecorder.Register(manager)
	if err != nil {
		return err
	}
	err = statusrecorder.Register(manager)
	if err != nil {
		return err
	}
	err = bindingrecorder.Register(manager)
	if err != nil {
		return err
	}
	err = endpointslicerecorder.Register(manager)
	if err != nil {
		return err
	}
	err = ownerreferencerecorder.Register(manager)
	if err != nil {
		return err
	}
	err = containerstatusrecorder.Register(manager)
	if err != nil {
		return err
	}
	err = snegrecorder.Register(manager)
	if err != nil {
		return err
	}
	err = noderecorder.Register(manager)
	if err != nil {
		return err
	}
	err = manager.Register(inspectionServer)
	if err != nil {
		return err
	}
	return nil
}
