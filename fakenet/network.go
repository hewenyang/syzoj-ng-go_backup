package fakenet

import (
	"context"
	"errors"
	"net"
	"sync"
)

type Network struct {
	m        sync.Mutex
	listener map[string]*listener
}

type listener struct {
	n    *Network
	ch   chan net.Conn
	done chan struct{}
	once sync.Once
	addr string
}

var ErrConnRefused = errors.New("fakenet: connection refused")
var ErrAddrInUse = errors.New("fakenet: address already in use")
var ErrClosed = errors.New("fakenet: socket closed")

var Base = NewNetwork()

func NewNetwork() *Network {
	n := &Network{}
	n.listener = make(map[string]*listener)
	return n
}

func (n *Network) Listen(addr string) (net.Listener, error) {
	n.m.Lock()
	_, found := n.listener[addr]
	if found {
		n.m.Unlock()
		return nil, ErrAddrInUse
	}
	l := &listener{}
	l.n = n
	l.ch = make(chan net.Conn)
	l.done = make(chan struct{})
	l.addr = addr
	n.listener[addr] = l
	n.m.Unlock()
	return l, nil
}

func (n *Network) DialContext(ctx context.Context, addr string) (net.Conn, error) {
	n.m.Lock()
	l, found := n.listener[addr]
	if !found {
		n.m.Unlock()
		return nil, ErrConnRefused
	}
	n.m.Unlock()
	a, b := net.Pipe()
	select {
	case <-ctx.Done():
		// note that currently net.Pipe does not leak goroutines
		return nil, ctx.Err()
	case <-l.done:
		return nil, ErrConnRefused
	case l.ch <- a:
		return b, nil
	}
}

func (l *listener) Accept() (net.Conn, error) {
	select {
	case <-l.done:
		return nil, ErrClosed
	case c := <-l.ch:
		return c, nil
	}
}

func (l *listener) Close() error {
	l.once.Do(func() {
		close(l.done)
		l.n.m.Lock()
		delete(l.n.listener, l.addr)
		l.n.m.Unlock()
	})
	return nil
}

func (l *listener) Addr() net.Addr {
	return fakeaddr(l.addr)
}

type fakeaddr string

func (f fakeaddr) Network() string {
	return "fakenet"
}

func (f fakeaddr) String() string {
	return string(f)
}
