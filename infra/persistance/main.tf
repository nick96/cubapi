provider "digitalocean" {
  token = var.token
}

resource "digitalocean_database_cluster" "clust" {
  name       = var.cluster_name
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = var.region
  node_count = var.node_count
  tags       = var.tags
}

resource "digitalocean_database_db" "db" {
  name       = var.db_name
  cluster_id = digitalocean_database_cluster.clust.id
}

resource "digitalocean_database_firewall" "fw" {
  cluster_id = digitalocean_database_cluster.clust.id

  rule {
    type  = "tag"
    value = var.allowed_tag
  }
}

resource "digitalocean_volume" "letsencrypt-vol" {
  region = var.region
  name   = var.volume_name
  size   = var.volume_size
}
