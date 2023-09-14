package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"scrapey/iocopy"
	"strconv"
)

func main() {
	proxy := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodConnect {
			client, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				w.WriteHeader(502)
				fmt.Println(err)
				return
			}
			defer client.Close()
			proxy, err := net.Dial("tcp", "127.0.0.1:5000")
			if err != nil {
				w.WriteHeader(400)
				fmt.Println(err)
				return
			}
			err = sendHost(proxy, req.URL.Host)
			if err != nil {
				w.WriteHeader(500)
				fmt.Println(err)
				return
			}
			client.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
			iocopy.Between(client, proxy)
		} else {
			w.WriteHeader(400)
			msg := "Only HTTPS supported."
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
			fmt.Fprintln(w, msg)
		}
	}
	fmt.Println("starting server on 5001")
	log.Fatal(http.ListenAndServe(":5001", http.HandlerFunc(proxy)))
}

func sendHost(conn net.Conn, host string) error {
	conn.Write([]byte(host))
	return nil
}
