/*
  Simulates 10 simultaneous requests to illustrate the effectiveness of
  multiplexing requests to a reverse HTTP tunnel, which maintains
  a single TCP connection, using HTTP/2.
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q := "http://localhost:8888/foo"
			fmt.Println("GET", "http://localhost:8888/foo")
			resp, err := http.DefaultClient.Get(q)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				b, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(resp.Status, string(b))
			} else {
				io.Copy(os.Stdout, resp.Body)
			}
		}()
	}
	wg.Wait()
}
