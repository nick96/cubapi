// Create the server infra structure
// - Droplet
// - Firewall

provider "digitalocean" {
  token = var.token
}

data "digitalocean_droplet_snapshot" "web_snapshot" {
  name   = var.image
  region = var.region
}

resource "digitalocean_droplet" "web" {
  image  = data.digitalocean_droplet_snapshot.web_snapshot.id
  name   = "web"
  region = var.region
  size   = "s-1vcpu-1gb"
  tags   = var.tags
}


resource "digitalocean_floating_ip_assignment" "ip_assign" {
  ip_address = var.public_ip
  droplet_id = digitalocean_droplet.web.id
}

resource "digitalocean_firewall" "web" {
  name        = "only-22-80-and-443"
  droplet_ids = [digitalocean_droplet.web.id]

  // HTTP
  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  // HTTPS
  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  // Ping
  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  // DNS
  outbound_rule {
    protocol              = "tcp"
    port_range            = "53"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  // DNS
  outbound_rule {
    protocol              = "udp"
    port_range            = "53"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  // Ping
  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}