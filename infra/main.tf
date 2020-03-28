
// Floating IP + DNS
module "network" {
  source = "./network"

  token  = var.do_token
  domain = "bagheera.nickspain.dev"
  region = var.do_region
}

// Droplet + firewall
module "server" {
  source = "./server"

  token     = var.do_token
  region    = var.do_region
  public_ip = module.network.public_ip
  tags      = ["bagheera"]
  image     = "bagheera-dev"
}

// DB (cluster and DB) + firewall
module "db" {
  source = "./db"

  token        = var.do_token
  cluster_name = "bagheera-db-clust"
  db_name      = "bagheera_db"
  tags         = ["bagheera"]
  allowed_tag  = "bagheera"
  node_count   = 1
  region       = var.do_region
}

provider "digitalocean" {
  token   = var.do_token
  version = "1.15.1"
}


data "digitalocean_project" "bagheera" {
  name = "bagheera"
}

resource "digitalocean_project_resources" "bagheera_resources" {
  project   = data.digitalocean_project.bagheera.id
  resources = concat(module.network.resources, module.server.resources, module.db.resources)
}
