# Kafka configuration

## Pre requisites
- A [Kafka](https://kafka.apache.org) Cluster visible to AsyncAPI Event-Gateway.

## Configuration
The Kafka proxy can be configured (mostly) by reading config from `servers` on an AsyncAPI document.  
However, some advanced configuration is only available through environment variables.

### From AsyncAPI doc
By reading config from `servers`, the Kafka proxy is configured by default as it follows:

- Uses each server `url` field as a new remote broker URL, using the same remote broker port as matching local port.
  - In case `EVENTGATEWAY_KAFKA_PROXY_BROKER_FROM_SERVER` environment variable is set, only that specific server is configured.
  - Local port can also be specified via `x-eventgateway-listener`, for example: ` x-eventgateway-listener: 28002`.
- In case `x-eventgateway-dial-mapping` extension is present, is used in the same form as the `EVENTGATEWAY_KAFKA_PROXY_BROKERS_DIAL_MAPPING` environment variable.


### From environment variables
| Environment variable                                | Type    | Description                                                                                                                                                                                                                                          | Default | Required                          | examples                                                                                                                                                                                                     |
|-----------------------------------------------------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|-----------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| EVENTGATEWAY_KAFKA_PROXY_BROKER_FROM_SERVER         | string  | When configuring from an AsyncAPI doc, this allows the user to only configure one server instead of all                                                                                                                                              | -       | No                                | `name-of-server1`, `server-test`                                                                                                                                                                             |
| EVENTGATEWAY_KAFKA_PROXY_BROKERS_MAPPING            | string  | Configure the mapping between remote broker address (the address published by the broker) and desired local address. Format is `remotehost:remoteport,localhost:localport`. Multiple values can be configured by using pipe separation (`\|`)        | -       | Yes when no AsyncAPI doc provided | `test.mykafkacluster.org:8092,localhost:28002`, `test.mykafkacluster.org:8092,localhost:28002|test2.mykafkacluster.org:8092,localhost:28003`                                                                 |
| EVENTGATEWAY_KAFKA_PROXY_BROKERS_DIAL_MAPPING       | string  | Configure the mapping between published remote broker address and the address the proxy will use when forwarding requests. Format is `remotehost:remoteport,localhost:localport`. Multiple values can be configured by using pipe separation (`\|`)  | -       | No                                | `test.mykafkacluster.org:8092,private-test.mykafkacluster.org:8092`, `test.mykafkacluster.org:8092,private-test.mykafkacluster.org:8092|test2.mykafkacluster.org:8092,private-test2.mykafkacluster.org:8092` |
| EVENTGATEWAY_KAFKA_PROXY_MESSAGE_VALIDATION_ENABLED | boolean | Enable or disable validation of Kafka messages                                                                                                                                                                                                       | `true`  | No                                | `true`, `false`                                                                                                                                                                                              |
| EVENTGATEWAY_KAFKA_PROXY_EXTRA_FLAGS                | string  | Advanced configuration. Configure any flag from [here](https://github.com/grepplabs/kafka-proxy/blob/4f3b89fbaecb3eb82426f5dcff5f76188ea9a9dc/cmd/kafka-proxy/server.go#L85-L195). Multiple values can be configured by using pipe separation (`\|`) | -       | No                                | `tls-enable=true|tls-client-cert-file=/opt/var/service.cert|tls-client-key-file=/opt/var/service.key`                                                                                                        |

