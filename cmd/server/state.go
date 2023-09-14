package main

import (
	"log"
	"net"
	"net/http"
)

type ServerState struct {
	proxies []net.Conn
}

func (s *ServerState) getProxy() (net.Conn, bool) {
	proxyCount := len(s.proxies)
	if proxyCount == 0 {
		return nil, false
	}
	proxy := s.proxies[proxyCount-1]
	s.proxies = s.proxies[:proxyCount-1]
	return proxy, true
}

func (s *ServerState) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		proxy, ok := s.getProxy()
		if !ok {
			handleNoProxies(w)
			return
		}
		handleConnectMethod(proxy, w, req)
	} else {
		handleNonConnectMethod(w)
	}
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
			return err
		}
		s.proxies = append(s.proxies, conn)
	}
}
