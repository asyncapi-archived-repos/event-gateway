<h5 align="center">
  <br>
  <a href="https://www.asyncapi.org"><img src="https://github.com/asyncapi/parser-nodejs/raw/master/assets/logo.png" alt="AsyncAPI logo" width="200"></a>
  <br>
  AsyncAPI Event Gateway
</h5>
<p align="center">
  <em>The Event Gateway solution by excellence</em>
</p>

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

# Goals

## 1. Performance-first while ensuring message delivery
The **Event Gateway** it's a **stateless** solution that ensures messages are delivered as fast as possible using a minimal resource footprint.
Delivering messages as a top priority means no data loss should happen.

## 2. Transparent usage.
No change in the user's code is needed. The service acts as a proxy between the client and the final broker(s). 
Messages infer the protocol based on the shape of the input network packet.

## 3. Fully configurable.
The service is entirely configurable, and the user can specify the settings for all protocols as well. For example, consumers' and producers' settings.

## 4. API-first
The service provides an API for uploading AsyncAPI specs, allowing the user to update their message validation, among others, very quickly. 
It could even be an automated task whenever you update your specs.

## 5. Extensible
The **Event Gateway** can extend its functionality via middlewares written by the community.
A catalog of middlewares made by the community is also available.