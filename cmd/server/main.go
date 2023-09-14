package main

import (
	"log"
	"net/http"
	"scrapey/utils"
)

func main() {
	serverPort := utils.EnvOrDefault("SERVER_PORT", "5001")
	proxyPort := utils.EnvOrDefault("PROXY_PORT", "5000")
	server := &ServerState{}

	go func() {
		err := server.listenForProxies(":" + proxyPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("server listening on port:", serverPort)
	err := http.ListenAndServe(":"+serverPort, server)
	if err != nil {
		log.Fatal(err)
	}
}
