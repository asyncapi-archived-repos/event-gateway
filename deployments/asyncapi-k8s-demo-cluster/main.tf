terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_kubernetes_cluster" "asyncapi-k8s-cluster" {
  name   = var.cluster_name
  region = var.region
  version = var.cluster_version

  node_pool {
    name       = var.cluster_node_pool_name
    size       = var.cluster_node_size
    auto_scale = true
    min_nodes  = var.cluster_node_count_min
    max_nodes  = var.cluster_node_count_max
  }
}
