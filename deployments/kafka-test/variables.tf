variable "aiven_api_token" {
  type = string
}

variable "aiven_project_name" {
  type = string
}

variable "aiven_service_name" {
  type = string
}

variable "aiven_service_plan" {
  type = string
}

variable "aiven_kafka_cloud_name" {
  type = string
}

variable "aiven_kafka_topics" {
  type = set(string)
}

variable "aiven_kafka_topics_partitions" {
  type = number
}

variable "aiven_kafka_topics_replication" {
  type = number
}

variable "aiven_kafka_topics_retention_ms" {
  type = number
}

variable "aiven_kafka_username" {
  type = string
}