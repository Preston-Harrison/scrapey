package main

import (
	"log"
	"net/http"
	"os"
)

func or(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

func main() {
	serverPort := or(os.Getenv("SERVER_PORT"), "5001")
	proxyPort := or(os.Getenv("PROXY_PORT"), "5000")

	server := &ServerState{}
	log.Println("server listening on port", serverPort)

	go func() {
		err := server.listenForProxies(":" + proxyPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := http.ListenAndServe(":"+serverPort, server)
	if err != nil {
		log.Fatal(err)
	}
}
