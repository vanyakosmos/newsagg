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

# storage

resource "ovh_cloud_project_user" "s3_user" {
  service_name = var.ovh_project_id
  description  = "S3 Operator for newsagg objects storage"
  role_name    = "objectstore_operator"
}

resource "ovh_cloud_project_user_s3_credential" "s3_user_creds" {
  service_name = var.ovh_project_id
  user_id      = ovh_cloud_project_user.s3_user.id
}

resource "ovh_cloud_project_storage" "newsagg_bucket" {
  service_name = var.ovh_project_id
  region_name  = var.ovh_region
  name         = "newsagg"
  owner_id     = ovh_cloud_project_user.s3_user.id
  limit        = 0
}

resource "ovh_cloud_project_user_s3_policy" "s3_user_policy" {
  service_name = var.ovh_project_id
  user_id      = ovh_cloud_project_user.s3_user.id
  policy = jsonencode({
    "Statement" : [{
      "Action" : ["s3:*"],
      "Effect" : "Allow",
      "Resource" : ["arn:aws:s3:::${ovh_cloud_project_storage.newsagg_bucket.name}", "arn:aws:s3:::${ovh_cloud_project_storage.newsagg_bucket.name}/*"],
      "Sid" : "FullAccess"
    }]
  })
}

# VPS

resource "ovh_cloud_project_ssh_key" "newsagg_admin" {
  service_name = var.ovh_project_id
  name         = "newsagg-admin-key"
  public_key   = file("~/.ssh/id_rsa.pub")
}
