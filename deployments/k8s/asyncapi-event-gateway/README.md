# asyncapi-event-gateway

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0-alpha](https://img.shields.io/badge/AppVersion-0.1.0--alpha-informational?style=flat-square)

Helm chart that installs the AsyncAPI Event-Gateway - https://github.com/asyncapi/event-gateway

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| asyncapiFileContent | string | `""` | AsyncAPI file. Should be the content of the file as plain text. Set either the AsyncAPI file content with '--set-file asyncapi-event-gateway.asyncapiFileContent=<filename>' or a URL to a valid spec file through '--set asyncapi-event-gateway.env.EVENTGATEWAY_ASYNC_API_DOC=<url>' The reason is that files from parent charts can't be accessed from subcharts. See https://github.com/helm/helm/pull/10077. |
| asyncapiFileMountPath | string | `"/app/asyncapi.yaml"` | The path where the AsyncAPI file will be located within the pod. In case you want to load the file from a URL, unset this value. |
| autoscaling.enabled | bool | `true` |  |
| autoscaling.maxReplicas | int | `4` |  |
| autoscaling.minReplicas | int | `1` |  |
| autoscaling.targetCPUUtilizationPercentage | int | `80` |  |
| env.EVENTGATEWAY_ASYNC_API_DOC | string | `"/app/asyncapi.yaml"` | This is where the asyncapi.yaml file is mounted when `--set-file asyncapi-event-gateway.asyncapiFileContent=event-gateway-demo/event-gateway-demo.asyncapi.yaml`. |
| env.EVENTGATEWAY_DEBUG | string | `"true"` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `asyncapi/event-gateway"` |  |
| image.tag | string | `"latest"` |  |
| imagePullSecrets | list | `[]` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| ports | object | `{"brokers":[],"healthcheck":80,"websocket":5000}` | Event-Gateway opened ports. |
| ports.brokers | list | `[]` | Specify ports for all possible brokers (both boostrap and discovered). |
| ports.healthcheck | int | `80` | Used as simple healthcheck. Called by K8s Deployment LivenessProbe and ReadinessProbe. |
| ports.websocket | int | `5000` | The websocket where the Event-Gateway API will be available. |
| replicaCount | int | `2` |  |
| resources | object | `{}` |  |
| secret | object | `{}` | Create a secret if needed. Useful for creating certificates for connecting to clusters such as Kafka. Combine it with a mounted volume, then you can load the certificates from a known path. |
| securityContext | object | `{}` |  |
| service.annotations | object | `{}` |  |
| service.type | string | `"ClusterIP"` |  |
| tolerations | list | `[]` |  |
| volumeMounts | list | `[]` |     secret:      secretName: "syncapi-event-gateway-foobar-kafka-certificates" |
| volumes | list | `[]` | Mount your volume. Especially useful for mounting secrets as explained above. |

