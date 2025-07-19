terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.44.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.6.0"
    }
  }
  required_version = ">= 1.12.2"
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
}

provider "github" {
  owner = var.github_owner
  token = var.github_token
}

resource "google_artifact_registry_repository" "docker" {
  location      = var.gcp_region
  repository_id = "docker"
  format        = "DOCKER"
  vulnerability_scanning_config {
    enablement_config = "DISABLED"
  }
}
