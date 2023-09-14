package main

import (
	"fmt"
	"net"
	"net/http"
	"scrapey/iocopy"
	"strconv"
)

func handleNonConnectMethod(w http.ResponseWriter) {
	w.WriteHeader(400)
	msg := "Only HTTPS supported."
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	fmt.Fprintln(w, msg)
}

func handleNoProxies(w http.ResponseWriter) {
	w.WriteHeader(500)
	msg := "No proxies available"
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	fmt.Fprintln(w, msg)
}

func handleConnectMethod(proxy net.Conn, w http.ResponseWriter, req *http.Request) {
	client, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		w.WriteHeader(502)
		fmt.Println(err)
		return
	}
	err = sendHost(proxy, req.URL.Host)
	if err != nil {
		client.Close()
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}
	client.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	iocopy.Between(client, proxy)
}

func sendHost(conn net.Conn, host string) error {
	conn.Write([]byte(host))
	return nil
}
