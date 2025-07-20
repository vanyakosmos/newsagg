terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.44.0"
    }
  }
  required_version = ">= 1.12.2"
  backend "gcs" {
    bucket = "newsagg-tf-state"
    prefix = "deploy"
  }
}

locals {
  gcp_project    = "newsagg-466216"
  gcp_region     = "europe-central2"
  job_name       = "newsagg"
  target_channel = "@newnewsagg"
}

variable "image_tag" {
  type = string
}

provider "google" {
  project = local.gcp_project
  region  = local.gcp_region
}

data "terraform_remote_state" "infra" {
  backend = "gcs"
  config = {
    bucket = "newsagg-tf-state"
    prefix = "infra"
  }
}

resource "google_cloud_run_v2_job" "newsagg" {
  name                = local.job_name
  location            = local.gcp_region
  deletion_protection = false

  template {
    parallelism = 1
    template {
      service_account = data.terraform_remote_state.infra.outputs.application_sa
      max_retries     = 0
      timeout         = "120s"
      containers {
        name  = "newsagg"
        image = "${data.terraform_remote_state.infra.outputs.image_name}:${var.image_tag}"
        resources {
          limits = {
            "memory" = "512Mi"
            "cpu" : "1"
          }
        }
        env {
          name  = "TELEGRAM_BOT_TOKEN"
          value = data.terraform_remote_state.infra.outputs.telegram_bot_token
        }
        env {
          name  = "TARGET_CHANNEL"
          value = local.target_channel
        }
        env {
          name  = "SENTRY_DSN"
          value = data.terraform_remote_state.infra.outputs.sentry_dsn
        }
        env {
          name  = "BUCKET_ENDPOINT"
          value = "storage.googleapis.com"
        }
        env {
          name  = "BUCKET_ACCESS_KEY"
          value = data.terraform_remote_state.infra.outputs.bucket_access_key
        }
        env {
          name  = "BUCKET_SECRET_KEY"
          value = data.terraform_remote_state.infra.outputs.bucket_secret_key
        }
        env {
          name  = "BUCKET_REGION"
          value = "auto"
        }
        env {
          name  = "BUCKET_NAME"
          value = data.terraform_remote_state.infra.outputs.bucket_name
        }
      }
    }
  }
}

resource "google_cloud_scheduler_job" "job" {
  name     = "${local.job_name}-scheduler"
  schedule = "*/30 * * * *"
  region   = local.gcp_region

  retry_config {
    retry_count = 0
  }

  http_target {
    http_method = "POST"
    uri         = "https://${local.gcp_region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${local.gcp_project}/jobs/${local.job_name}:run"
    oauth_token {
      service_account_email = data.terraform_remote_state.infra.outputs.cloudscheduler_sa
    }
  }
}
