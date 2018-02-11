# PTTH

Package `ptth` is a proof-of-concept implementation of the [ReverseHTTP](http://reversehttp.net/) idea. Originally inspired by this Tweet from the inimitable [@KelseyHightower](https://twitter.com/kelseyhightower):

> "HTTP tunneling has the potential to dramatically simplify service discovery, end to end security, and configuration of API gateways and reverse proxies." - https://twitter.com/kelseyhightower/status/950371704504598529


## What's in this project?

This project includes a router implementation and a way to easily establish reverse HTTP tunnels to the router from backend services. It also provides several examples demonstrating how to use the package.

The router accepts HTTP traffic and acts as a load-balancer and reverse-proxy to the backend services. Backend services dial the router on a TCP port, which is separate from the port the router is listening for HTTP traffic on. As such, no ingress configuration is required for the backend services.

## What's not in this project?

Anything production ready ;)

This is just a proof of concept. As such, _you will find no security here_. Specifically not implemented is any TLS handshaking or [h2c](https://http2.github.io/http2-spec/#iana-h2c) upgrades.

However, if the tunneling port is only exposed to services running within a private network (such as an Amazon VPC), then it could be used safely...in theory. Use at your own risk.

## Using the examples

Start a reverse proxy that accepts ReverseHTTP tunnel connections:

```sh
$ go run _exampes/router.go
2018/02/11 17:19:00 Listening for reverse tunnels on tcp://localhost:8887
2018/02/11 17:19:00 Listening for HTTP traffic on http://localhost:8888
```

Start a backend web service that handles `GET /foo` requests:

```sh
$ go run _examples/backend.go
2018/02/11 17:19:43 Dialing tcp://localhost:8887 and serving HTTP/2 traffic
```

Execute some HTTP traffic to the router:

```sh
$ go run _examples/client.go
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 > GET http://localhost:8888/foo
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
2018/02/11 17:25:04 < 200 OK "bar"
```