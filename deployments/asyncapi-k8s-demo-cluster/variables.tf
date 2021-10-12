# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}

variable "region" {
  type = string
  default = "nyc1" # Regions list: `doctl kubernetes options regions`
}

variable "cluster_name" {
  type = string
  default = "asyncapi-demo"
}

variable "cluster_version" {
  type = string
  default = "1.21.5-do.0" # Versions list: `doctl kubernetes options versions`
}

variable "cluster_node_pool_name" {
  type = string
  default = "asyncapi-demo-pool"
}

variable "cluster_node_size" {
  type = string
  default = "s-2vcpu-4gb" # Sizes list: `doctl kubernetes options sizes`
}

variable "cluster_node_count_min" {
  type = number
  default = 3
}

variable "cluster_node_count_max" {
  type = number
  default = 5
}