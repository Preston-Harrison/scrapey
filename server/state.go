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

func (s *ServerState) popProxy(channelKey string) (net.Conn, bool) {
	proxyList, ok := s.proxies[channelKey]
	if !ok {
		return nil, false
	}
	return proxyList.pop()
}

func (s *ServerState) addProxy(channelKey string, proxy net.Conn) {
	proxyList, ok := s.proxies[channelKey]
	if !ok {
		s.m.Lock()
		defer s.m.Unlock()
		s.proxies[channelKey] = &ProxyList{&sync.Mutex{}, []net.Conn{proxy}}
	} else {
		proxyList.push(proxy)
	}
}
