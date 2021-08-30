# AsyncAPI Event-Gateway config reference
The Event-Gateway is configured through a combination of environment variables and one AsyncAPI document.  
AsyncAPI documents can be used for configuring some proxies (from Servers, Channels, etc.). However, advanced configuration is only possible through environment variables.  
Configuration for Event-Gateway is done via environment variables as well.

## Global
| Environment variable        | Type    | Description                                                    | Default | Required | examples                                                                                                    |
|-----------------------------|---------|----------------------------------------------------------------|---------|----------|-------------------------------------------------------------------------------------------------------------|
| EVENTGATEWAY_DEBUG          | boolean | Enable or disable debug logs                                              | `false` | No       | `true`, `false`                                                                                             |
| EVENTGATEWAY_ASYNC_API_DOC  | string  | Path or URL to a valid AsyncAPI doc (v2.0.0 is only supported) | `false` | No       | `/var/opt/streetlights.yml`, `https://github.com/asyncapi/spec/blob/master/examples/streetlights-kafka.yml` |
| EVENTGATEWAY_WS_SERVER_PORT | integer | Port for the Websocket server. Used for debugging events       | `5000`  | No       | `5000`, `9000`                                                                                              |


### Protocol specific
- [Kafka](kafka.md)