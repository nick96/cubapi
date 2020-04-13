variable "do_token" {
  type = string
}

variable "do_region" {
  type    = string
  default = "sgp1"
}

variable "letsencrypt_volume_name" {
  type    = string
  default = "letsencrypt"
}

variable "database_tag" {
  type    = string
  default = "bagheera-db"
}
