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
  backend "gcs" {
    bucket = "newsagg-tf-state"
    prefix = "infra"
  }
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
}

provider "github" {
  owner = var.github_owner
  token = var.github_token
}
