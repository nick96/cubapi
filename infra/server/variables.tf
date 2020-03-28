// Token for accessing digital ocean
variable "token" {
  type = string
}

// Region in which to create everything
variable "region" {
  type    = string
  default = "sgp1"
}

// Public IP address, this is what the DNS record points to
variable "public_ip" {
  type = string
}

// Tags to apply to the droplet
variable "tags" {
  type    = list(string)
  default = []
}

// Image to build the droplet from
variable "image" {
  type = string
}
