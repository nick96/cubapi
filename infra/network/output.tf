output "public_ip" {
  value = digitalocean_floating_ip.public_ip.ip_address
}

output "resources" {
  value = [digitalocean_domain.domain.urn]
}
