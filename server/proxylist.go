package main

import (
	"net"
	"sync"
)

type ProxyList struct {
	m     *sync.Mutex
	conns []net.Conn
}

func (p *ProxyList) pop() (net.Conn, bool) {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.conns) == 0 {
		return nil, false
	}
	proxy := p.conns[len(p.conns)-1]
	p.conns = p.conns[:len(p.conns)-1]
	return proxy, true
}

func (p *ProxyList) push(c net.Conn) {
	p.m.Lock()
	defer p.m.Unlock()
	p.conns = append([]net.Conn{c}, p.conns...)
}