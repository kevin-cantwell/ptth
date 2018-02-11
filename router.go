package ptth

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

// Router manages a pool of reverse HTTP tunnels and routes HTTP
// traffic to them in a randomized pattern.
type Router struct {
	mu   sync.Mutex
	pool []*Tunnel
}

// ListenAndAcceptTunnels listens for TCP connections on addr
// and adds them to a pool.
func (router *Router) ListenAndAcceptTunnels(addr string) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer ln.Close()
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				log.Println("Error accepting tunnel conn:", err)
				time.Sleep(time.Second)
				continue
			}

			router.mu.Lock()
			router.pool = append(router.pool, NewTunnel(conn))
			router.mu.Unlock()

			log.Println("Tunnel added:", conn.RemoteAddr().String())
		}
	}()
}

// ServeHTTP proxies the request to a reverse HTTP tunnel. The
// tunnel itself is chosen randomly from a pool. If no healthy tunnel
// is available, an error of http.StatusServiceUnavailable will be
// served.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tunnel := router.getTunnel()
	if tunnel == nil {
		http.Error(w, "Server at capacity", http.StatusServiceUnavailable)
		return
	}
	log.Println("Proxying", r.Method+" "+r.URL.Path, "to", tunnel.conn.RemoteAddr().String())
	tunnel.ServeHTTP(w, r)
}

func (router *Router) getTunnel() *Tunnel {
	router.mu.Lock()
	defer router.mu.Unlock()

	// Keep trying to find a healthy tunnel until the pool exhausted
	for {
		if len(router.pool) < 1 {
			return nil
		}
		i := rand.Intn(len(router.pool))
		tunnel := router.pool[i]
		if tunnel.Closed() {
			router.pool = append(router.pool[:i], router.pool[i+1:]...)
			log.Println("Tunnel removed:", tunnel.conn.RemoteAddr().String())
			continue
		}
		return tunnel
	}
}
