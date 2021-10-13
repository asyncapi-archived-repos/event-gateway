# Kafka configuration

## Pre requisites
- A [Kafka](https://kafka.apache.org) Cluster visible to AsyncAPI Event-Gateway.

## Configuration
The Kafka proxy is mostly configured by setting some properties on the `server` object.

| Server property         | Type   | Description                                                                                                                                                                                                                                                | Default                        | Required | examples                                                                                                                                      |
|-------------------------|--------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| x-eventgateway-listener | string | Configure the mapping between remote broker address (the address set in the broker server `URL` field) and desired local address. Format is `remotehost:remoteport,localhost:localport`. Multiple values can be configured by using pipe separation (`\|`) | `0.0.0.0:<remote-server-port>` | No       | `test.mykafkacluster.org:8092,localhost:28002`, `test.mykafkacluster.org:8092,localhost:28002\|test2.mykafkacluster.org:8092,localhost:28003` |
| x-eventgateway-dial-mapping | string | Override a broker address with another address the proxy will use when connecting. Format is `remotehost:remotehost,new-remotehost:new-remoteport`. Multiple values can be configured by using pipe separation (`\|`) | -                              | No       | `0.0.0.0:8092,test.myeventgateway.org:8092`, `test.mykafkacluster.org:8092,mykafkacluster.org:8092,\|`test.mykafkacluster.org:8093,mykafkacluster.org:8093`          |

#### Example
```yaml
# ...
servers:
  test:
    url: broker.mybrokers.org:9092
    protocol: kafka
    x-eventgateway-listener: 28002 # optional. 0.0.0.0:9092 will be used instead if missing.
    x-eventgateway-dial-mapping: '0.0.0.0:28002,test.myeventgateway.org:8092' # optional. 
# ...
```

## Advanced configuration
Some advanced configuration is only available through environment variables.

| Environment variable                                | Type    | Description                                                                                                                                                                                                                                          | Default   | Required | examples                                                                                                |
|-----------------------------------------------------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|----------|---------------------------------------------------------------------------------------------------------|
| EVENTGATEWAY_KAFKA_PROXY_ADDRESS                    | string  | Address for this proxy. Clients will use this address as host when connecting to the brokers through this proxy, so it should be reachable by your clients. Most probably a domain name.                                                             | `0.0.0.0` | No       | `event-gateway-demo.asyncapi.org`                                                                       |
| EVENTGATEWAY_KAFKA_PROXY_BROKER_FROM_SERVER         | string  | When set, only the specified server will be considered instead of all servers.                                                                                                                                                                       | -         | No       | `name-of-server1`, `server-test`                                                                        |
| EVENTGATEWAY_KAFKA_PROXY_MESSAGE_VALIDATION_ENABLED | boolean | Enable or disable validation of Kafka messages                                                                                                                                                                                                       | `true`    | No       | `true`, `false`                                                                                         |
| EVENTGATEWAY_KAFKA_PROXY_EXTRA_FLAGS                | string  | Advanced configuration. Configure any flag from [here](https://github.com/grepplabs/kafka-proxy/blob/4f3b89fbaecb3eb82426f5dcff5f76188ea9a9dc/cmd/kafka-proxy/server.go#L85-L195). Multiple values can be configured by using pipe separation (`\|`) | -         | No       | `tls-enable=true\|tls-client-cert-file=/opt/var/service.cert\|tls-client-key-file=/opt/var/service.key` |