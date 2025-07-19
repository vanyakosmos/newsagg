resource "google_storage_bucket" "newsagg" {
  name                     = "newsagg"
  location                 = var.gcp_region
  public_access_prevention = "enforced"
}

resource "google_storage_hmac_key" "newsagg" {
  service_account_email = google_service_account.application.email
}
