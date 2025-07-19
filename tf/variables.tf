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
  type = string
}
variable "github_repo" {
  type    = string
  default = "newsagg"
}

variable "telegram_bot_token" {
  type = string
}
