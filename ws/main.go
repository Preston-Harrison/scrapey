package main

import (
	"fmt"
	"net"
	"scrapey/iocopy"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Println("listening on 5000")
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go tcpForward(conn)
	}
}

func tcpForward(conn net.Conn) {
	defer conn.Close()
	host, err := readHost(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("host", host)

	hostConn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		return
	}

	iocopy.Between(conn, hostConn)
}

func readHost(client net.Conn) (string, error) {
	h := make([]byte, 1024)
	i, err := client.Read(h)
	if err != nil {
		return "", err
	}
	return string(h[:i]), nil
}