/*
  Simulates 10 simultaneous requests to illustrate the effectiveness of
  multiplexing requests to a reverse HTTP tunnel, which maintains
  a single TCP connection, using HTTP/2.
*/
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q := "http://localhost:8888/foo"
			log.Println("> GET", "http://localhost:8888/foo")
			resp, err := http.DefaultClient.Get(q)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			log.Printf("< %s %q\n", resp.Status, string(b))
		}()
	}
	wg.Wait()
}
