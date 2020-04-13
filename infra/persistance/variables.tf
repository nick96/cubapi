// Digital ocean token
variable "token" {
  type = string
}

// Name of the cluster to create
variable "cluster_name" {
  type = string
}

// Region to place the cluster in
variable "region" {
  type    = string
  default = "sgp1"
}

// Number of nodes to spin up in the cluster
variable "node_count" {
  type    = number
  default = 1
}

// Tags to apply to the cluster
variable "tags" {
  type    = list(string)
  default = []
}

// Tag allowed to access the DB cluster
variable "allowed_tag" {
  type = string
}

// Name of the database to create
variable "db_name" {
  type = string
}

// Name of the volume to create
variable "volume_name" {
  type = string
}

// Size of the volume to create
variable "volume_size" {
  type = number
}
