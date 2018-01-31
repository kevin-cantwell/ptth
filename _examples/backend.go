package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hkwi/h2c"
)

type ListenerConn struct {
	conn net.Conn
}

func (l *ListenerConn) Accept() (net.Conn, error) {
	return l.conn, nil
}

func (l *ListenerConn) Close() error {
	return l.conn.Close()
}

func (l *ListenerConn) Addr() net.Addr {
	return l.conn.LocalAddr()
}

func main() {
	raddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling:", r.URL.Path)
		w.Write([]byte(r.URL.Path + "\n"))
	})

	s := http.Server{
		Handler: &h2c.Server{},
	}
	if err := s.Serve(&ListenerConn{conn: conn}); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
	log.Println("disconnected")
}
