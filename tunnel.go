package ptth

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/net/http2"
)

type Tunnel struct {
	proxy *httputil.ReverseProxy
	conn  *netConn
}

func NewTunnel(conn net.Conn) *Tunnel {
	nc := &netConn{Conn: conn}

	rURL, _ := url.Parse("http://" + conn.RemoteAddr().String())
	proxy := httputil.NewSingleHostReverseProxy(rURL) // The rURL value doesn't actually matter
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

func (t *Tunnel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.proxy.ServeHTTP(w, r)
}

func (t *Tunnel) Closed() bool {
	return t.conn.closed
}
