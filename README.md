# PTTH

Package `ptth` is a proof-of-concept implementation of the [ReverseHTTP](http://reversehttp.net/) idea. Originally inspired by this Tweet from the inimitable [@KelseyHightower](https://twitter.com/kelseyhightower):

https://twitter.com/kelseyhightower/status/950371704504598529

## What's in this project?

This project includes a router implementation and a way to easily establish reverse HTTP tunnels to the router from backend services. It also provides several examples demonstrating how to use the package.

The router accepts HTTP traffic and acts as a load-balancer and reverse-proxy to the backend services. Backend services dial the router on a TCP port, which is separate from the port the router is listening for HTTP traffic on. As such, no ingress configuration is required for the backend services.

## What's not in this project?

Anything production ready ;)

This is just a proof of concept. As such, _you will find no security here_. Specifically not implemented is any TLS handshaking or [h2c](https://http2.github.io/http2-spec/#iana-h2c) upgrades.

However, if the tunneling port is only exposed to services running within a private network (such as an Amazon VPC), then it could be used safely...in theory. Use at your own risk.