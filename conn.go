package ptth

import (
	"net"
)

var _ net.Conn = &netConn{} // Compiler to trick to ensure netConn always implements net.Conn

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
