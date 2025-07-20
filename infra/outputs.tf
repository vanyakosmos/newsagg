# used for reading from deploy

output "application_sa" {
  value = google_service_account.application.email
}

output "cloudscheduler_sa" {
  value = google_service_account.cloudscheduler.email
}

output "image_name" {
  value = local.image_name
}

output "telegram_bot_token" {
  value     = var.telegram_bot_token
  sensitive = true
}

output "sentry_dsn" {
  value     = var.sentry_dsn
  sensitive = true
}

output "bucket_access_key" {
  value     = google_storage_hmac_key.newsagg.access_id
  sensitive = true
}

output "bucket_secret_key" {
  value     = google_storage_hmac_key.newsagg.secret
  sensitive = true
}

output "bucket_name" {
  value = google_storage_bucket.newsagg.name
}
