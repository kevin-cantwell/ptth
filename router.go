package ptth

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

var (
	_ net.Conn     = &netConn{} // Compiler trick to ensure *netConn implements net.Conn
	_ http.Handler = &Tunnel{}  // Compiler trick to ensure that *Tunnel implements http.Handler
	_ http.Handler = &Router{}  // Compiler trick to ensure that *Router implements http.Handler
)

// netConn is a simple wrapper for net.Conn that keeps track of whether or not the
// connection is closed. We do this because the underlying connection is difficult
// to extract from the http and http2 packages. One could argue this is a dirty hack.
type netConn struct {
	net.Conn // Embedded net.Conn ensures net.Conn implementation
	closed   bool
}

func (c *netConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if err != nil {
		c.closed = true
	}
	return n, err
}

func (c *netConn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	if err != nil {
		c.closed = true
	}
	return n, err
}

// Tunnel multiplexes HTTP requests over HTTP/2 using a reverse proxy configured
// with a single instance of net.Conn. Implements http.Handler
type Tunnel struct {
	proxy *httputil.ReverseProxy
	conn  *netConn
}

// Provides a new instance of *Tunnel that i
func NewTunnel(conn net.Conn) *Tunnel {
	nc := &netConn{Conn: conn}

	// The rURL value doesn't actually matter, as we are not actually dialing to anything
	rURL, _ := url.Parse("http://" + conn.RemoteAddr().String())
	proxy := httputil.NewSingleHostReverseProxy(rURL)
	proxy.Transport = &http2.Transport{
		DialTLS: func(netw, addr string, _ *tls.Config) (net.Conn, error) {
			// HTTP/2 protocol normally requires a TLS handshake. This works
			// around that by using an already established connection. This
			// also avoids the usual requirement of performing an h2c upgrade
			// when not using TLS.
			return nc, nil
		},
		// Routed requests may use the http scheme if we specify this config.
		AllowHTTP: true,
	}

	return &Tunnel{
		proxy: proxy,
		conn:  nc,
	}
}

// ServeHTTP multiplexes requests on a single TCP connection by using
// HTTP/2.
func (t *Tunnel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.proxy.ServeHTTP(w, r)
}

// Closed indicates if there have been any read or write errors on the
// underlying connection. Clients may use this information to decide to
// evict the tunnel from a pool of tunnels.
func (t *Tunnel) Closed() bool {
	return t.conn.closed
}

// Router manages a pool of reverse HTTP tunnels and routes HTTP
// traffic to them in a randomized pattern.
type Router struct {
	mu   sync.Mutex
	pool []*Tunnel
}

// ListenAndAcceptTunnels listens for TCP connections on addr
// and adds them to a pool of reverse HTTP tunnels.
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

// ServeHTTP proxies the request to a single reverse HTTP tunnel. The
// tunnel itself is chosen randomly from a pool. If no healthy tunnel
// is available, an error of http.StatusServiceUnavailable will be
// served.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tunnel := router.GetTunnel()
	if tunnel == nil {
		http.Error(w, "Server at capacity", http.StatusServiceUnavailable)
		return
	}
	log.Println("Proxying", r.Method+" "+r.URL.Path, "to", tunnel.conn.RemoteAddr().String())
	tunnel.ServeHTTP(w, r)
}

// GetTunnel provides a random tunnel from a pool of tunnels. If no
// tunnels exist, the return value will be nil. Any non-nil *Tunnel
// returned is guaranteed to have had zero read or write errors. This
// does not guarantee that the underlying net.Conn is healthy.
func (router *Router) GetTunnel() *Tunnel {
	router.mu.Lock()
	defer router.mu.Unlock()

	// Keep trying to find a healthy tunnel until the pool exhausted and
	// actively remove unhealthy tunnels from the pool.
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
