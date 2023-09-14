package main

import (
	"log"
	"net"
	"scrapey/iocopy"
	"scrapey/utils"
)

func main() {
	server := utils.EnvOrDefault("SERVER_HOST", "127.0.0.1:5000")
	channelKey := utils.EnvOrPanic("CHANNEL_KEY")
	for {
		log.Println("waiting for new host")
		conn, err := net.Dial("tcp", server)
		if err != nil {
			panic(err)
		}
		attachToChannel(conn, channelKey)
		host, err := waitForHost(conn)
		if err != nil {
			log.Println("failed to wait for host:", err)
			continue
		}
		log.Println("proxying to host:", host)
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

func attachToChannel(conn net.Conn, c string) {
	conn.Write([]byte(c))
}