package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {
    proxyURL, err := url.Parse("http://localhost:5001")
    if err != nil {
        panic(err)
    }
    httpClient := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }
    response, err := httpClient.Get("https://example.com/")
	if err != nil {
		panic(err)
	}
    log.Println("Response status:", response.Status)
}
