package ptth

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

type Router struct {
	backends map[string]chan *net.TCPConn
}

func NewRouter() *Router {
	return &Router{
		backends: map[string]chan *net.TCPConn{},
	}
}

func (r *Router) ListenAndProxyHTTP(port string, h http.Handler) {
	if h == nil {
		h = http.DefaultServeMux
	}
	routingHandler := func(resp http.ResponseWriter, req *http.Request) {
		ch, ok := r.backends[req.Host]
		if !ok {
			http.Error(resp, http.StatusText(404), 404)
			return
		}
		conn := <-ch
	}
	http.ListenAndServe(":"+port, h)
}

func (r *Router) ListenForTCPBackends(port string) {
	serverAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		log.Fatal("Unable to listen on "+serverAddr.String()+": ", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println("Error accepting tcp conn:", err)
			continue
		}

		go func(conn *net.TCPConn) {
			addr := conn.RemoteAddr().String()
			log.Printf("Client connected: %s\n", addr)
			defer log.Printf("Client disconnected: %s\n", addr)
			defer conn.Close()

			host, err := r.handleHandshake(conn)
			if err != nil {
				log.Printf("Handshake error from %s: %v\n", addr, err)
				return
			}
			if err := r.proxyHost(host, conn); err != nil {
				log.Printf("Proxy error from %s: %v\n", addr, err)
				return
			}
		}(conn)
	}
}

func (r *Router) handleHandshake(conn *net.TCPConn) (string, error) {
	buf := bufio.NewReader(conn)
	return buf.ReadString('\n')
}

func (r *Router) proxyHost(host string, conn *net.TCPConn) error {
	ch, ok := r.backends[host]
	if !ok {
		ch = make(chan *net.TCPConn)
		r.backends[host] = ch
	}
	ch <- conn
}
