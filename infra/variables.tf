locals {
  image_name = "${github_actions_variable.docker_registry.value}/${var.gcp_project}/docker/newsagg"
}

variable "gcp_project" {
  type    = string
  default = "newsagg-466216"
}
variable "gcp_region" {
  type    = string
  default = "europe-central2"
}

variable "github_owner" {
  type    = string
  default = "vanyakosmos"
}
variable "github_token" {
  type      = string
  sensitive = true
}
variable "github_repo" {
  type    = string
  default = "newsagg"
}

variable "telegram_bot_token" {
  type      = string
  sensitive = true
}

variable "sentry_dsn" {
  type      = string
  sensitive = true
}
