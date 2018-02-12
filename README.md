# PTTH

Package `ptth` is a proof-of-concept implementation of the [ReverseHTTP](http://reversehttp.net/) idea. Originally inspired by this Tweet from the inimitable [@KelseyHightower](https://twitter.com/kelseyhightower):

> "HTTP tunneling has the potential to dramatically simplify service discovery, end to end security, and configuration of API gateways and reverse proxies." - https://twitter.com/kelseyhightower/status/950371704504598529

I wanted to know what an implementation of this would look like--if it could be done easily in Go--so I rolled my own. It was a teachable moment for me and I hope it helps others imagine what a ReverseHTTP solution might look like.


## What's in this project?

This project includes a router implementation and a package function that backend services may use to establish reverse HTTP tunnels to the router. It also provides a working example of each.

The router accepts HTTP traffic and acts as a load-balancer and reverse-proxy to the backend services. Backend services dial the router on a TCP port (which is separate from the HTTP port the router accepts user requests on) and establishes a reverse HTTP tunnel, handling all requests over HTTP/2. As such, no ingress configuration is required for the backend services. Nor is any service discovery needed beyond knowing the router address. This is the point of ReverseHTTP ;)

I should mention that [backplane.io](https://backplane.io) already offers production-grade ReverseHTTP as a service. Their solution uses a host agent that performs the tunneling magic and acts as a reverse proxy on localhost. My implementation is an in-process solution that is for illustration purposes only.

## What's not in this project?

Anything production ready ;)

This is just a proof of concept. As such, _you will find no security here_. Specifically not implemented is any TLS handshaking or [h2c](https://http2.github.io/http2-spec/#iana-h2c) upgrades.

However, if the tunneling port is only exposed to services running within a private network (such as an Amazon VPC), then it could be used safely...in theory. Use at your own risk.

## Using the examples

Start a HTTP server that routes requests to ReverseHTTP tunnels:

```sh
$ go run _exampes/router.go
2018/02/11 17:19:00 Listening for reverse tunnels on tcp://localhost:8887
2018/02/11 17:19:00 Listening for HTTP traffic on http://localhost:8888
```

Start a backend web service that tunnels into the router and handles `GET /foo` requests:

```sh
$ go run _examples/backend.go
2018/02/11 17:19:43 Dialing tcp://localhost:8887 and serving HTTP/2 traffic
```

Execute some user requests on the router:

```sh
$ go run _examples/user.go
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