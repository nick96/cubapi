output "resources" {
  value = [digitalocean_database_cluster.clust.urn, digitalocean_volume.letsencrypt-vol.urn]
}

output "volume_id" {
  value = digitalocean_volume.letsencrypt-vol.id
}

output "database_port" {
  value = digitalocean_database_cluster.clust.port
}
