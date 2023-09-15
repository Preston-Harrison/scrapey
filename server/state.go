package main

import (
	"net"
	"sync"
)

type ServerState struct {
	m *sync.Mutex
	// Map of channel key to proxy list.
	proxies map[string]*ProxyList
}

func NewServerState() *ServerState {
	return &ServerState{&sync.Mutex{}, make(map[string]*ProxyList)}
}

func (s *ServerState) getProxy(authToken string) (net.Conn, bool) {
	proxyList, ok := s.proxies[authToken]
	if !ok {
		return nil, false
	}
	return proxyList.pop()
}

func (s *ServerState) addProxy(authToken string, proxy net.Conn) {
	proxyList, ok := s.proxies[authToken]
	if !ok {
		s.m.Lock()
		defer s.m.Unlock()
		s.proxies[authToken] = &ProxyList{&sync.Mutex{}, []net.Conn{proxy}}
	} else {
		proxyList.push(proxy)
	}
}
