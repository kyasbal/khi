package task

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/env"

const EnvGcpIamTokenVariableName = GCPPrefix + "environment/iam-token"

var EnvGcpIamTokenTask = env.EnvironmentVariableProducer(EnvGcpIamTokenVariableName, "IAM_TOKEN", "")

const EnvFixedProjectIdVariableName = GCPPrefix + "environment/fixed-project-id"

var EnvFixedProjectIdTask = env.EnvironmentVariableProducer(EnvFixedProjectIdVariableName, "KHI_FIXED_PROJECT_ID", "")

const EnvQuotaProjectIdVariableName = GCPPrefix + "environment/quota-project-id"

var EnvQuotaProjectIdTask = env.EnvironmentVariableProducer(EnvQuotaProjectIdVariableName, "KHI_QUOTA_PROJECT_ID", "")
