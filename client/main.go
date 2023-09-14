package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				User: url.UserPassword("user", "1234"),
				Host: "127.0.0.1:5001",
			}),
		},
	}
	response, err := httpClient.Get("https://example.com/")
	if err != nil {
		panic(err)
	}
	log.Println("Response status:", response.Status)
}
