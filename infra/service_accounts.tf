resource "google_service_account" "github" {
  account_id   = "github"
  display_name = "github"
  description  = "Github CI/CD"
  project      = var.gcp_project
}
resource "google_storage_bucket_iam_member" "gitlab" {
  member = "serviceAccount:${google_service_account.github.email}"
  bucket = "newsagg-tf-state"
  role   = "roles/storage.objectAdmin"
}
resource "google_project_iam_member" "github_artifactregistry_writer" {
  member  = "serviceAccount:${google_service_account.github.email}"
  project = var.gcp_project
  role    = "roles/artifactregistry.writer"
}
resource "google_project_iam_member" "github_service_account_user" {
  member  = "serviceAccount:${google_service_account.github.email}"
  project = var.gcp_project
  role    = "roles/iam.serviceAccountUser"
}
resource "google_project_iam_member" "github_run_developer" {
  member  = "serviceAccount:${google_service_account.github.email}"
  project = "newsagg-466216"
  role    = "roles/run.developer"
}
resource "google_project_iam_member" "github_cloudscheduler_admin" {
  member  = "serviceAccount:${google_service_account.github.email}"
  project = var.gcp_project
  role    = "roles/cloudscheduler.admin"
}
resource "google_service_account_key" "github" {
  service_account_id = google_service_account.github.account_id
  keepers = {
    key_version = 1
  }
}


resource "google_service_account" "cloudscheduler" {
  account_id   = "cloudscheduler"
  display_name = "cloudscheduler"
  description  = "Schedules cloud run functions/jobs"
  project      = var.gcp_project
}
resource "google_project_iam_member" "cloudscheduler_run_job_executor" {
  member  = "serviceAccount:${google_service_account.cloudscheduler.email}"
  project = var.gcp_project
  role    = "roles/run.jobsExecutor"
}


resource "google_service_account" "application" {
  account_id   = "application"
  display_name = "application"
  description  = "Used inside the app to access stuff"
  project      = var.gcp_project
}
resource "google_project_iam_member" "application_storage_object_admin" {
  member  = "serviceAccount:${google_service_account.application.email}"
  project = var.gcp_project
  role    = "roles/storage.objectAdmin"
}
