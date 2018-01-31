package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

type http2ClientPool struct {
	conn      net.Conn
	transport *http2.Transport
}

func (pool *http2ClientPool) GetClientConn(req *http.Request, addr string) (*http2.ClientConn, error) {
	return nil, nil
}

func (pool *http2ClientPool) MarkDead(*http2.ClientConn) {

}

func main() {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, r.Proto+"\n")
		})
		// log.Fatal(http.ListenAndServe(":6666", &h2c.Server{}))

		laddr, err := net.ResolveTCPAddr("tcp", ":7777")
		if err != nil {
			log.Fatal(err)
		}
		ln, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			log.Fatal("Unable to listen on "+laddr.String()+": ", err)
		}
		defer ln.Close()

		s := http2.Server{}
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				log.Fatalln("Error accepting tcp conn:", err)
			}
			go s.ServeConn(conn, &http2.ServeConnOpts{
			// Handler: &h2c.Server{},
			})
		}

	}()

	transport := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(netw, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
		// ConnPool:
	}
	// transport.ConnPool = &http2ClientPool{conn: }
	client := http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", "http://localhost:7777", nil)
	if err != nil {
		log.Fatalln("Unabled to create request:", err)
	}
	// req.Header.Set("Host", "localhost")
	// req.Header.Set("Connection", "Upgrade, HTTP2-Settings")
	// req.Header.Set("Upgrade", "h2c")
	// req.Header.Set("HTTP2-Settings", "AAAABAAAAAAA")

	for {
		func() {
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
		}()
	}
}
