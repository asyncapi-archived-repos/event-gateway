terraform {
  required_providers {
    aiven = {
      source  = "aiven/aiven"
      version = "2.1.19"
    }
  }
}

provider "aiven" {
  api_token = var.aiven_api_token
}

data "aiven_project" "asyncapi-project" {
  project = var.aiven_project_name
}

# Kafka service
resource "aiven_kafka" "asyncapi-event-gateway-kafka-service" {
  project                 = data.aiven_project.asyncapi-project.project
  cloud_name              = var.aiven_kafka_cloud_name
  plan                    = var.aiven_service_plan
  service_name            = var.aiven_service_name
  maintenance_window_dow  = "friday"
  maintenance_window_time = "07:00:00"
  termination_protection = true
  kafka_user_config {
    kafka_version = "2.7"
  }
}

# Topic for Kafka
resource "aiven_kafka_topic" "asyncapi-event-gateway-kafka-topic" {
  project      = data.aiven_project.asyncapi-project.project
  service_name = aiven_kafka.asyncapi-event-gateway-kafka-service.service_name
  topic_name   = var.aiven_kafka_topic_name
  partitions   = var.aiven_kafka_topic_partitions
  replication  = var.aiven_kafka_topic_replication
  config {
    retention_ms = var.aiven_kafka_topic_retention_ms
  }
}

# User for Kafka
resource "aiven_service_user" "asyncapi-event-gateway-kafka-user" {
  project      = data.aiven_project.asyncapi-project.project
  service_name = aiven_kafka.asyncapi-event-gateway-kafka-service.service_name
  username     = var.aiven_kafka_username
}

# ACL for Kafka
resource "aiven_kafka_acl" "asyncapi-event-gateway-read-acl" {
  project      = data.aiven_project.asyncapi-project.project
  service_name = aiven_kafka.asyncapi-event-gateway-kafka-service.service_name
  username     = aiven_service_user.asyncapi-event-gateway-kafka-user.username
  permission   = "read"
  topic        = aiven_kafka_topic.asyncapi-event-gateway-kafka-topic.topic_name
}

# ACL for Kafka
resource "aiven_kafka_acl" "asyncapi-event-gateway-write-acl" {
  project      = data.aiven_project.asyncapi-project.project
  service_name = aiven_kafka.asyncapi-event-gateway-kafka-service.service_name
  username     = aiven_service_user.asyncapi-event-gateway-kafka-user.username
  permission   = "write"
  topic        = aiven_kafka_topic.asyncapi-event-gateway-kafka-topic.topic_name
}