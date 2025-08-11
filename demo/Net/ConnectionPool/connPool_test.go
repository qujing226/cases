package test

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type ConnectionPool struct {
	mu          *sync.RWMutex
	connections map[string]*PooledConnection
	maxConns    int
	idleTimeout time.Duration
}

type PooledConnection struct {
	transport   *NetTransport
	lastUsed    time.Time
	inUse       int32
	maxRequests int32
}

func (p *ConnectionPool) GetConnection(address string) (*NetTransport, error) {
	p.mu.RLock()
	if conn, exists := p.connections[address]; exists {
		if atomic.LoadInt32(&conn.inUse) < conn.maxRequests {
			atomic.AddInt32(&conn.inUse, 1)
			p.mu.RUnlock()
			return conn.transport, nil
		}
	}
	p.mu.RUnlock()
	return p.createConnection(address)
}

func (p *ConnectionPool) createConnection(address string) (*NetTransport, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.connections) >= p.maxConns {
		return nil, errors.New("too many connections")
	}
	transport := &NetTransport{}
}

type NetTransport struct {
}

type MultiplexTransport struct {
	pool         *ConnectionPool
	loadBalancer LoadBalancer
}
