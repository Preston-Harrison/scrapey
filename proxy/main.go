package main

import (
	"log"
	"net"
	"os"
	"scrapey/iocopy"
)

func or(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

func main() {
	server := or(os.Getenv("SERVER"), "127.0.0.1:5000")
	for {
		log.Println("waiting for new host")
		conn, err := net.Dial("tcp", server)
		if err != nil {
			panic(err)
		}
		host, err := waitForHost(conn)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("proxying to host", host)
		go tcpForward(conn, host)
	}
}

func waitForHost(conn net.Conn) (string, error) {
	h := make([]byte, 1024)
	i, err := conn.Read(h)
	if err != nil {
		return "", err
	}
	return string(h[:i]), nil
}

func tcpForward(proxyConn net.Conn, host string) error {
	hostConn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	iocopy.Between(proxyConn, hostConn)
	return nil
}