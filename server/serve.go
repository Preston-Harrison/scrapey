package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func (s *ServerState) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		_, password, err := parseProxyAuth(req)
		if err != nil {
			handleBadAuthHeader(w, err)
			return
		}
		log.Println("received request for channel key:", password)
		proxy, ok := s.popProxy(password)
		if !ok {
			handleNoProxies(w)
			return
		}
		handleConnectMethod(proxy, w, req)
	} else {
		handleNonConnectMethod(w)
	}
}

// Right now channel key is just parsed from the proxy authorization password.
func parseProxyAuth(req *http.Request) (string, string, error) {
	authHeader := req.Header.Get("Proxy-Authorization")
	if authHeader == "" {
		return "", "", fmt.Errorf("Proxy-Authorization header not found")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", fmt.Errorf("invalid Proxy-Authorization header format")
	}

	credentials, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", err
	}

	credentialsParts := strings.SplitN(string(credentials), ":", 2)
	if len(credentialsParts) != 2 {
		return "", "", fmt.Errorf("invalid credentials format")
	}

	username := credentialsParts[0]
	password := credentialsParts[1]

	return username, password, nil
}

func (s *ServerState) listenForProxies(host string) error {
	log.Println("listening for proxies on host:", host)
	l, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("cannot accept tcp connection:", err)
			conn.Close()
			continue
		}
		channelKey := make([]byte, 1024)
		i, err := conn.Read(channelKey)
		if err != nil {
			log.Println("failed to read channel key", err)
			conn.Close()
			continue
		}
		keyStr := string(channelKey[:i])
		log.Printf("adding new proxy for key %s\n", keyStr)
		s.addProxy(keyStr, conn)
	}
}
