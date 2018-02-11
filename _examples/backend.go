/*
  This is an example of a backend web service that serves traffic over a reverse
  HTTP tunnel instead of listening on a port. The service is responsible for
  dialing a single TCP connection to a router and all requests are multiplexed
  over the connection using HTTP/2.
*/
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/kevin-cantwell/ptth"
)

func main() {
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling:", r.URL.Path)
		time.Sleep(time.Millisecond) // Arbitrary delay to help illustrate HTTP/2 multiplexing
		w.Write([]byte("bar"))
	})

	log.Println("Dialing tcp://localhost:8887 and serving HTTP/2 traffic")
	if err := ptth.DialRouterAndServe(":8887", nil); err != nil {
		log.Println(err)
	}
}
