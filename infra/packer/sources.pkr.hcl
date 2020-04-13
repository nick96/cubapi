source "digitalocean" "base_image" {
  ssh_username  = "root"
  api_token     = var.do_token
  image         = "ubuntu-18-04-x64"
  region        = "sgp1"
  size          = "s-1vcpu-1gb"
  snapshot_name = "bagheera-${var.version}"
}
