
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

  token                 = var.do_token
  region                = var.do_region
  public_ip             = module.network.public_ip
  tags                  = ["bagheera"]
  image                 = "bagheera-dev"
  letsencrypt_volume_id = module.persistance.volume_id
  database_port         = module.persistance.database_port
}

// DB (cluster and DB) + firewall
module "persistance" {
  source = "./persistance"

  token        = var.do_token
  cluster_name = "bagheera-db-clust"
  db_name      = "bagheera_db"
  tags         = [var.database_tag]
  allowed_tag  = "bagheera"
  node_count   = 1
  region       = var.do_region
  volume_name  = var.letsencrypt_volume_name
  volume_size  = 1
}

provider "digitalocean" {
  token   = var.do_token
  version = "1.15.1"
}


data "digitalocean_project" "bagheera" {
  name = "bagheera"
}

resource "digitalocean_project_resources" "bagheera_resources" {
  project = data.digitalocean_project.bagheera.id
  resources = concat(
    module.network.resources,
    module.server.resources,
    module.persistance.resources,
  )
}
