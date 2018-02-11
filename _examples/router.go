/*
  This is an example of a reverse proxy that serves all HTTP requests
  to a pool of reverse HTTP tunnels. It listens on two ports: One public
  port that accepts incoming HTTP requests and one private port that
  accepts TCP connections and adds them to a pool of long-lived reverse HTTP
  tunnels. Because each tunnel maintains a single TCP connection to a backend
  service, incoming HTTP requests are proxied using HTTP/2.
*/
package main

import (
	"log"
	"net/http"

	"github.com/kevin-cantwell/ptth"
)

func main() {
	var router ptth.Router

	log.Println("Listening for reverse tunnels on tcp://localhost:8887")
	go router.ListenAndAcceptTunnels(":8887")

	log.Println("Listening for HTTP traffic on http://localhost:8888")
	http.ListenAndServe(":8888", &router)
}
