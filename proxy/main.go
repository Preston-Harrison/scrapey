package main

import (
	"log"
	"net"
	"scrapey/handshake"
	"scrapey/iocopy"
	"scrapey/utils"
)

func main() {
	server := utils.EnvOrDefault("SERVER_HOST", "127.0.0.1:5000")
	authToken := utils.EnvOrPanic("AUTH_TOKEN")
	for {
		log.Println("waiting for new host")
		conn, err := net.Dial("tcp", server)
		if err != nil {
			panic(err)
		}
		err = handshake.SendAuthToken(conn, authToken)
		if err != nil {
			panic(err)
		}
		target, err := handshake.ReceiveAndDialHost(conn)
		if err != nil {
			log.Println("failed to receive host", err)
			continue
		}
		go iocopy.Between(conn, target)
	}
}
