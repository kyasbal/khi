provider "google" {
  project = "kubernetes-history-inspector"
}

data "google_iam_policy" "plx_service_account" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = [
      "user:ikakeru@prod.google.com",
    ]
  }

  binding {
    role = "roles/iam.serviceAccountTokenCreator"

    members = [
      "user:plx-security@prod.google.com",
    ]
  }
}

resource "google_service_account" "plx_service_account" {
  account_id   = "plx-service-account"
  display_name = "Service account for fetching BigQuery data in PLX"
}

resource "google_service_account_iam_policy" "plx_service_account" {
  service_account_id = google_service_account.plx_service_account.name
  policy_data        = data.google_iam_policy.plx_service_account.policy_data
}

resource "google_project_iam_member" "plx_service_account_data_viewer" {
  role    = "roles/bigquery.dataEditor"
  project = "kubernetes-history-inspector"
  member  = "serviceAccount:${google_service_account.plx_service_account.email}"
}

resource "google_project_iam_member" "plx_service_account_job_user" {
  role    = "roles/bigquery.jobUser"
  project = "kubernetes-history-inspector"
  member  = "serviceAccount:${google_service_account.plx_service_account.email}"
}

resource "google_project_iam_member" "khi_testing_service_account_data_viewer" {
  role    = "roles/bigquery.dataEditor"
  project = "kubernetes-history-inspector"
  member  = "serviceAccount:composer-env-account@khi-testing.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "khi_testing_service_account_job_user" {
  role    = "roles/bigquery.jobUser"
  project = "kubernetes-history-inspector"
  member  = "serviceAccount:composer-env-account@khi-testing.iam.gserviceaccount.com"
}


resource "google_project_iam_binding" "project-ownwer" {
  project = "kubernetes-history-inspector"
  role    = "roles/owner"
  members = [
    "user:ikakeru@google.com"
  ]
}

resource "google_project_iam_binding" "project-editor" {
  project = "kubernetes-history-inspector"
  role    = "roles/editor"
  members = [
    "user:ryunosukes@google.com",
    "user:takeie@google.com",
    "serviceAccount:22938676367@cloudservices.gserviceaccount.com"
  ]
}

resource "google_project_iam_binding" "image-reader" {
  project = "kubernetes-history-inspector"
  role    = "roles/artifactregistry.reader"
  members = [
    "domain:google.com",
    "domain:premium-cloud-support.com",
    "serviceAccount:515136028509@cloudbuild.gserviceaccount.com"
  ]
}

resource "google_project_iam_binding" "logging-reader" {
  project = "kubernetes-history-inspector"
  role    = "roles/logging.privateLogViewer"
  members = [
    "serviceAccount:471017113462@cloudbuild.gserviceaccount.com",
    "domain:google.com",
    "domain:premium-cloud-support.com"
  ]
}

resource "google_project_iam_binding" "bq-reader" {
  project = "kubernetes-history-inspector"
  role    = "roles/logging.privateLogViewer"
  members = [
    "serviceAccount:471017113462@cloudbuild.gserviceaccount.com",
    "domain:google.com",
    "domain:premium-cloud-support.com"
  ]
}

resource "google_project_iam_binding" "build-job-invoker" {
  project = "kubernetes-history-inspector"
  role    = "roles/cloudbuild.builds.editor"
  members = [
    "domain:google.com",
    "domain:premium-cloud-support.com"
  ]
}

resource "google_cloud_run_service_iam_binding" "analytics-invoker" {
  location = "us-central1"
  project  = "kubernetes-history-inspector"
  service  = "khi-analytics"
  role     = "roles/run.invoker"
  members = [
    "domain:google.com",
    "domain:premium-cloud-support.com"
  ]
}
