output "_project_info" {
  value = data.ovh_cloud_project.newsagg
}

output "blob_storage_endpoint" {
  value = ovh_cloud_project_storage.newsagg_bucket.virtual_host
}

output "blob_access_key" {
  description = "the access key that have been created by the terraform script"
  value       = ovh_cloud_project_user_s3_credential.s3_user_creds.access_key_id
}

output "blob_secret_key" {
  description = "the secret key that have been created by the terraform script"
  value       = ovh_cloud_project_user_s3_credential.s3_user_creds.secret_access_key
  sensitive   = true
}
