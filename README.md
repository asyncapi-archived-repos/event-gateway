[![AsyncAPI Event Gateway](./assets/logo.png)](https://www.asyncapi.com)

<h5 align="center">
  <br>
  <a href="https://github.com/asyncapi/event-gateway/issues/new?assignees=&labels=use+case&template=use_case.md&title=%5BUSECASE%5D+">
    <img src="https://dummyimage.com/1000x80/0e9f6f/ffffff.png&text=We+are+looking+for+use+cases!+Please+share+yours+by+clicking+here" alt="Share your use case with us">
  </a>
  <br>
</h5>

> :warning: Still under development, it didn't reach v1.0.0 and therefore is not suitable for production use yet.

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-6-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

# Overview  

AsyncAPI **Event Gateway** (name is subject to change) is the Event Gateway solution by excellence.

Based on traditional API Gateways, it intercepts all incoming messages moving them into a pipeline of middlewares and handlers such as:

- Message validation
- Message manipulation
- Message aggregation
- Message filtering
- Authentication
- Throttling
- Routing
- Monitoring (including tracing)

It supports all the protocols AsyncAPI supports through [bindings](https://github.com/asyncapi/bindings).

This **Event Gateway** is also compatible with the HTTP protocol, natively or through an external provider like [Krakend.io](http://krakend.io).

![AsyncAPI Event Gateway big picture](https://user-images.githubusercontent.com/1083296/120669755-07323e00-c490-11eb-8844-a6292b516656.jpg)

## Goals

### 1. Performance-first while ensuring message delivery
The **Event Gateway** it's a **stateless** solution that ensures messages are delivered as fast as possible using a minimal resource footprint.
Delivering messages as a top priority means no data loss should happen.

### 2. Transparent usage.
No change in the user's code is needed. The service acts as a proxy between the client and the final broker(s). 
Messages infer the protocol based on the shape of the input network packet.

### 3. Fully configurable.
The service is entirely configurable, and the user can specify the settings for all protocols as well. For example, consumers' and producers' settings.

### 4. API-first
The service provides an API for uploading AsyncAPI specs, allowing the user to update their message validation, among others, very quickly. 
It could even be an automated task whenever you update your specs.

### 5. Extensible
The **Event Gateway** can extend its functionality via middlewares written by the community.
A catalog of middlewares made by the community is also available.

## Roadmap
The idea is to keep iterating and support all the protocols AsyncAPI supports through [bindings](https://github.com/asyncapi/bindings).  
However, we reduced the scope for the first versions, so we can give support to the most used protocols. 

For the first version, only [Kafka](https://kafka.apache.org) protocol will be supported. 

## Demo
A demo of the **Event Gateway** has been deployed and it is completely available to the users.
It is a very limited demo environment, but it is a good starting point for the users to familiarize on the concept.

Note that as per today, only the [Kafka](https://kafka.apache.org) protocol is supported. More info at [docs/config/kafka.md](docs/config/kafka.md).
The AsyncAPI file used for the demo is available [here](deployments/k8s/event-gateway-demo/event-gateway-demo.asyncapi.yaml). You can also [open it with Studio](https://studio.asyncapi.com/?url=https://raw.githubusercontent.com/asyncapi/event-gateway/master/deployments/k8s/event-gateway-demo/event-gateway-demo.asyncapi.yaml).

This video introduces this demo application:  
[![Quick introduction to what AsyncAPI Event-Gateway is capable of in 2021](https://user-images.githubusercontent.com/1083296/146788873-f21fca6b-a07e-4cac-841d-0f59f48065bd.png)](https://www.youtube.com/watch?v=ht8mOf3wsFw)

Main things you can do now:

- **Produce** Kafka messages to the Kafka broker `event-gateway-demo.asyncapi.com:20472` and topic `event-gateway-demo`. The expected message payload (`lightMeasured`) is:
  ```yaml
  type: object
  properties:
    lumens:
      type: integer
      minimum: 0
      description: Light intensity measured in lumens.
    sentAt:
      type: string
      format: date-time
      description: Date and time when the message was sent.
  ```
  You can use [kcat](https://github.com/edenhill/kcat) as the Kafka producer (headers are totally optional):
  ```bash
  echo '{"lumens": 100}' | kcat -H test=true -H origin=readme -H time=(date) -b event-gateway-demo.asyncapi.com:20472 -t event-gateway-demo -P
  ```                
- **Consume** Kafka messages from the Kafka broker `event-gateway-demo.asyncapi.com:20472` and topic `event-gateway-demo`.
  You can use [kcat](https://github.com/edenhill/kcat) as the Kafka consumer (headers are totally optional):
  ```bash
  kcat -b event-gateway-demo.asyncapi.com:20472 -t event-gateway-demo -C
  ```     
- **Watch** for the messages via the websocket endpoint at `ws://event-gateway-demo.asyncapi.com:5000/ws`. You can use [Websocat](https://github.com/vi/websocat) to connect to the websocket:
  ```bash
  websocat -v ws://event-gateway-demo.asyncapi.com:5000/ws
  ```
  This websocket endpoint consumes from another Kafka topic where the original produced messages are forwarded. 
  
  In the case of producing messages with an invalid payload (payloads that don't validate against the schema above), an extra metadata field (Kafka header under the hood) named `_asyncapi_eg_validation_error` will be included in the message. For example:
  ```json
  "_asyncapi_eg_validation_error": "{\"ts\":\"2021-12-20T11:33:26.583143572Z\",\"errors\":[\"lumens: Invalid type. Expected: integer, given: boolean\"]}",
  ```

## Getting Started

### Install from Docker
TBD

### Install from pre-compiled binaries
TBD

### Install from source
This project is built with [Go](https://golang.org/), and it uses [Go Modules](https://golang.org/ref/mod) for managing dependencies.  
The Minimum required version of Go is set in [go.mod](go.mod) file.

1. Clone this repository.
2. Run `make build`. The binary will be placed at `bin/out/event-gateway`.

### Configuration
Please refer to [config reference](./docs/config/README.md).

## Contributing
Read [CONTRIBUTING](https://github.com/asyncapi/.github/blob/master/CONTRIBUTING.md) guide.

## Contributors ✨
Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/smoya"><img src="https://avatars.githubusercontent.com/u/1083296?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Sergio Moya</b></sub></a><br /><a href="#question-smoya" title="Answering Questions">💬</a> <a href="https://github.com/asyncapi/event-gateway/issues?q=author%3Asmoya" title="Bug reports">🐛</a> <a href="https://github.com/asyncapi/event-gateway/commits?author=smoya" title="Code">💻</a> <a href="https://github.com/asyncapi/event-gateway/commits?author=smoya" title="Documentation">📖</a> <a href="#ideas-smoya" title="Ideas, Planning, & Feedback">🤔</a> <a href="#infra-smoya" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-smoya" title="Maintenance">🚧</a> <a href="#projectManagement-smoya" title="Project Management">📆</a> <a href="#research-smoya" title="Research">🔬</a> <a href="https://github.com/asyncapi/event-gateway/pulls?q=is%3Apr+reviewed-by%3Asmoya" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/asyncapi/event-gateway/commits?author=smoya" title="Tests">⚠️</a></td>
    <td align="center"><a href="http://www.fmvilas.com/"><img src="https://avatars.githubusercontent.com/u/242119?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Fran Méndez</b></sub></a><br /><a href="#ideas-fmvilas" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/asyncapi/event-gateway/pulls?q=is%3Apr+reviewed-by%3Afmvilas" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="https://github.com/magicmatatjahu"><img src="https://avatars.githubusercontent.com/u/20404945?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Maciej Urbańczyk</b></sub></a><br /><a href="https://github.com/asyncapi/event-gateway/pulls?q=is%3Apr+reviewed-by%3Amagicmatatjahu" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="https://dev.to/derberg"><img src="https://avatars.githubusercontent.com/u/6995927?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Lukasz Gornicki</b></sub></a><br /><a href="https://github.com/asyncapi/event-gateway/pulls?q=is%3Apr+reviewed-by%3Aderberg" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="http://polr.fr/me"><img src="https://avatars.githubusercontent.com/u/904193?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Paul B.</b></sub></a><br /><a href="https://github.com/asyncapi/event-gateway/pulls?q=is%3Apr+reviewed-by%3ApaulRbr" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="https://github.com/jonaslagoni"><img src="https://avatars.githubusercontent.com/u/13396189?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Jonas Lagoni</b></sub></a><br /><a href="https://github.com/asyncapi/event-gateway/pulls?q=is%3Apr+reviewed-by%3Ajonaslagoni" title="Reviewed Pull Requests">👀</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
