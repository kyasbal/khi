resource "google_project_iam_audit_config" "artifact-registry-audit" {
  project = "kubernetes-history-inspector"
  service = "artifactregistry.googleapis.com"
  audit_log_config {
    log_type = "ADMIN_READ"
  }
  audit_log_config {
    log_type = "DATA_READ"
  }
  audit_log_config {
    log_type = "DATA_WRITE"
  }
}

resource "google_bigquery_dataset" "audit-dataset" {
  dataset_id    = "khi_audit"
  friendly_name = "Audit log sink for KHI"
  description   = "A dataset stores artifact registry audit log data for monitoring KHI usage"
  location      = "US"
  access {
    role = "OWNER"
    user_by_email = "ikakeru@google.com"
  }
  access {
    role = "WRITER"
    user_by_email = "cloud-logs@system.gserviceaccount.com"
  }
  access {
    role   = "READER"
    domain = "google.com"
  }
}

resource "google_service_account" "service_account" {
  account_id   = "audit-bq-writer"
  display_name = "Cloud Logging Sink router Service account"
}

resource "google_logging_project_sink" "artifact-registry-audit-logging-sink"{
    name = "artifact-registry-audit-logging-sink"
    destination = "bigquery.googleapis.com/projects/kubernetes-history-inspector/datasets/${google_bigquery_dataset.audit-dataset.dataset_id}"
    filter="protoPayload.serviceName=\"artifactregistry.googleapis.com\""
}