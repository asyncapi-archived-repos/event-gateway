# event-gateway-demo

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0-alpha](https://img.shields.io/badge/AppVersion-0.1.0--alpha-informational?style=flat-square)

Helm chart that installs a demo of the AsyncAPI Event-Gateway.

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../asyncapi-event-gateway | asyncapi-event-gateway | 0.1.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| asyncapi-event-gateway.env.EVENTGATEWAY_ASYNC_API_DOC | string | `"/app/asyncapi.yaml"` | This is where the asyncapi.yaml file is mounted when `--set-file asyncapi-event-gateway.asyncapiFileContent=event-gateway-demo/event-gateway-demo.asyncapi.yaml`. |
| asyncapi-event-gateway.env.EVENTGATEWAY_DEBUG | string | `"true"` |  |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_ADDRESS | string | `"event-gateway-demo.asyncapi.org"` | his is the address that points to the DO load balancer in front of the Event-Gateway K8s service. |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_BROKER_FROM_SERVER | string | `"asyncapi-kafka-test"` | Only use `asyncapi-kafka-test` declared server. |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_EXTRA_FLAGS | string | `"dynamic-sequential-min-port=20473"` | Dynamic broker listeners will start on this port. Aiven cluster has, at least, 3 brokers (3 discovered brokers, one of those is known as the bootstrap server (20472)). |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_MESSAGE_VALIDATION_PUBLISH_TO_KAFKA_TOPIC | string | `"event-gateway-demo-validation"` | event-gateway-demo-validation is the topic where validation errors will be published to. The app reads from it and exposes those errors through the ws server. |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_TLS_CA_CHAIN_CERT_FILE | string | `"/etc/certs/ca"` | Value comes from `--set asyncapi-event-gateway.secret.data.cert=$(base64 ca-file-path)` |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_TLS_CLIENT_CERT_FILE | string | `"/etc/certs/cert"` | Value comes from `--set asyncapi-event-gateway.secret.data.cert=$(base64 cert-file-path)` |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_TLS_CLIENT_KEY_FILE | string | `"/etc/certs/key"` | Value comes from `--set asyncapi-event-gateway.secret.data.cert=$(base64 key-file-path)` |
| asyncapi-event-gateway.env.EVENTGATEWAY_KAFKA_PROXY_TLS_ENABLE | string | `"true"` |  |
| asyncapi-event-gateway.ports | object | `{"brokers":[20472,20473,20474,20475]}` | As dynamic-sequential-min-port is set to 20473, and Aiven has 3 brokers, we add those to the list apart from the kwnown bootstrap one (20472). |
| asyncapi-event-gateway.resources | object | `{}` | TODO Set this once we run some load testing within DO. |
| asyncapi-event-gateway.secret | object | `{"data":{},"name":"asyncapi-event-gateway-aiven-certificates"}` | cert, key and ca set via `--set asyncapi-event-gateway.secret.data.{cert|key|ca}=$(base64 filepath)` |
| asyncapi-event-gateway.secret.name | string | `"asyncapi-event-gateway-aiven-certificates"` | Same name as the configured in `asyncapi-event-gateway.volumes[secret-volume].secret.secretName`. |
| asyncapi-event-gateway.service.annotations."service.beta.kubernetes.io/do-loadbalancer-protocol" | string | `"tcp"` |  |
| asyncapi-event-gateway.service.annotations."service.beta.kubernetes.io/do-loadbalancer-size-slug" | string | `"lb-small"` |  |
| asyncapi-event-gateway.service.type | string | `"LoadBalancer"` | LoadBalancer type will tell Digital Ocean K8s to create a Network load balancer based on the annotations set above. |
| asyncapi-event-gateway.volumeMounts[0].mountPath | string | `"/etc/certs"` |  |
| asyncapi-event-gateway.volumeMounts[0].name | string | `"secret-volume"` |  |
| asyncapi-event-gateway.volumeMounts[0].readOnly | bool | `true` |  |
| asyncapi-event-gateway.volumes[0].name | string | `"secret-volume"` |  |
| asyncapi-event-gateway.volumes[0].secret.secretName | string | `"asyncapi-event-gateway-aiven-certificates"` | Same name as the configured in `asyncapi-event-gateway.secret.name`. |

