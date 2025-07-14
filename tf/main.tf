terraform {
  required_providers {
    ovh = {
      source  = "ovh/ovh"
      version = "2.5.0"
    }
  }
}

provider "ovh" {
  endpoint           = "ovh-eu"
  application_key    = var.ovh_application_key
  application_secret = var.ovh_application_secret
  consumer_key       = var.ovh_consumer_key
}

data "ovh_cloud_project" "newsagg" {
  service_name = var.ovh_project_id
}

resource "ovh_cloud_project_user" "user" {
  service_name = var.ovh_project_id
  description  = "Object Storage User"
  role_name    = "objectstore_operator"
}

resource "ovh_cloud_project_user_s3_credential" "admin" {
  service_name = var.ovh_project_id
  user_id      = ovh_cloud_project_user.user.id
}

resource "ovh_cloud_project_storage" "blob_storage" {
  service_name = var.ovh_project_id
  region_name  = var.ovh_region
  name         = "newsagg"
  owner_id     = ovh_cloud_project_user.user.id
}
