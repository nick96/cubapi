variable "do_token" {}
variable "do_region" {
  default = "sgp1"
}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "local" {
  name = "local"
}


resource "digitalocean_volume" "cubapi_db" {
  region                  = var.do_region
  name                    = "cubapi-db"
  size                    = 5
  initial_filesystem_type = "ext4"
  description             = "Cub API DB volume"
}

resource "digitalocean_droplet" "web" {
  image      = "ubuntu-18-04-x64"
  name       = "web"
  region     = var.do_region
  size       = "s-1vcpu-1gb"
  ssh_keys   = [data.digitalocean_ssh_key.local.id]
  volume_ids = [digitalocean_volume.cubapi_db.id]
}


resource "digitalocean_floating_ip" "web_ip" {
  droplet_id = digitalocean_droplet.web.id
  region     = digitalocean_droplet.web.region
}

output "floating_ip" {
  value = digitalocean_floating_ip.web_ip.ip_address
}
