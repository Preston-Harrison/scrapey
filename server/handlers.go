package main

import (
	"fmt"
	"net"
	"net/http"
	"scrapey/handshake"
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

func handleBadAuthHeader(w http.ResponseWriter, err error) {
	w.WriteHeader(403)
	msg := err.Error()
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
		fmt.Println("failed to hijack response writer:", err)
		return
	}
	err = handshake.SendHost(proxy, req.URL.Host)
	if err != nil {
		client.Close()
		w.Write([]byte("HTTP/1.0 500 Not OK\r\n\r\n"))
		fmt.Println("failed to send host to proxy:", err)
		return
	}
	// Tell client that it can start sending bytes.
	client.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	iocopy.Between(client, proxy)
}
