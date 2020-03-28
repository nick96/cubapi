// Create the network infrastructure
// - Floating IP
// - DNS records


provider "digitalocean" {
  token = var.token
}

resource "digitalocean_domain" "domain" {
  name = var.domain
}


resource "digitalocean_floating_ip" "public_ip" {
  region     = var.region
}

resource "digitalocean_record" "www" {
  domain = digitalocean_domain.domain.name
  type = "A"
  name = "www"
  value = digitalocean_floating_ip.public_ip.ip_address
}

resource "digitalocean_record" "self" {
  domain = digitalocean_domain.domain.name
  type = "A"
  name = "@"
  value = digitalocean_floating_ip.public_ip.ip_address
}


resource "digitalocean_record" "star" {
  domain = digitalocean_domain.domain.name
  type = "A"
  name = "*"
  value = digitalocean_floating_ip.public_ip.ip_address
}

