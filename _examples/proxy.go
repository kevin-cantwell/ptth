package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	var tunnel net.Conn

	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		ln, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			log.Fatal("Unable to listen on "+laddr.String()+": ", err)
		}
		defer ln.Close()

		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatalln("Error accepting tcp conn:", err)
		}

		req, err := http.NewRequest("GET", "http://"+conn.RemoteAddr().String()+"/http2upgrade", nil)
		if err != nil {
			log.Fatalln("Unabled to create request:", err)
		}
		req.Header.Set("Host", "localhost")
		req.Header.Set("Connection", "Upgrade, HTTP2-Settings")
		req.Header.Set("Upgrade", "h2c")
		req.Header.Set("HTTP2-Settings", "AAAABAAAAAAA")
		c := http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return conn, nil
				},
			},
		}
		log.Println("Doing...")
		resp, err := c.Do(req)
		log.Println("Done...")
		if err != nil {
			log.Fatalln("Error upgrading to h2c:", err)
		}
		defer resp.Body.Close()

		tunnel = conn
		log.Println("Established http2 tunnel")
	}()

	fakeURL, err := url.Parse("http://fake.com")
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(fakeURL)
	proxy.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if tunnel == nil {
				return nil, errors.New("dial: no tunnel exists")
			}
			return tunnel, nil
		},
	}
	http.ListenAndServe(os.Args[1], proxy)
}
