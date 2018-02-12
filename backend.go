package ptth

import (
	"errors"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

// CreateTunnelAndServe establishes a tcp connection to addr and
// serves incoming HTTP/2 requests. Only a single connection
// is used to multiplex requests. An error will be returned
// if any of the following occur: 1) A tcp connection cannot
// be established to the router; 2) The connection is
// closed by the remote host; or 3) The http2 server stops serving
// requests due to an internal error.
func CreateTunnelAndServe(addr string, handler http.Handler) error {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	s := http2.Server{}
	s.ServeConn(conn, &http2.ServeConnOpts{
		// If handler is nil, defaults to http.DefaultServeMux
		Handler: handler,
	})
	// if we've reached this point, the reason is internal to
	// the http2 package. Setting http2.VerboseLogs can help debug this.
	return errors.New("ptth: http2 server stopped serving")
}
