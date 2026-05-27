package main

import (
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

var ErrPoolExhausted = errors.New("连接池已满")

type ConnectionPool struct {
	activeConns map[string]*PooledConnection
	idleConns   chan *PooledConnection
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration

	factory func() (net.Conn, error)
	mu      sync.RWMutex
	closed  bool
}

type PooledConnection struct {
	conn       net.Conn
	transport  Transport
	lastUsed   time.Time
	inUse      bool
	clientAddr string
}

func (p *ConnectionPool) GetConnection(clientAddr string) (*PooledConnection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 1. 检查是否又该客户端的活跃链接
	if conn, exists := p.activeConns[clientAddr]; exists && conn.inUse {
		conn.inUse = true
		conn.lastUsed = time.Now()
		return conn, nil
	}
	// 2. 从空闲池获取连接
	select {
	case conn := <-p.idleConns:
		conn.inUse = true
		conn.lastUsed = time.Now()
		p.activeConns[clientAddr] = conn
		return conn, nil
	default:
		// 3. 创建新连接（未达到最大限制）
		if len(p.activeConns) < p.maxActive {
			return p.createNewConnection(clientAddr)
		}
		return nil, ErrPoolExhausted
	}

}
func (p *ConnectionPool) createNewConnection(clientAddr string) (*PooledConnection, error) {
	conn, err := p.factory()
	if err != nil {
		return nil, err

	}
	return &PooledConnection{
		conn:       conn,
		transport:  Transport{},
		lastUsed:   time.Now(),
		inUse:      true,
		clientAddr: clientAddr,
	}, nil
}

func (p *ConnectionPool) ReleaseConnection(conn *PooledConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn.inUse = false
	conn.lastUsed = time.Now()

	// 检查空闲池是否已满
	select {
	case p.idleConns <- conn:
		// 成功放入空闲池
	default:
		// 空闲池已满，关闭连接
		conn.conn.Close()
		delete(p.activeConns, conn.clientAddr)
	}
}

func (p *ConnectionPool) cleanupExpiredConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.mu.Lock()
			now := time.Now()

			// 清理超时的空闲连接
			for len(p.idleConns) > 0 {
				select {
				case conn := <-p.idleConns:
					if now.Sub(conn.lastUsed) > p.idleTimeout {
						conn.conn.Close()
						delete(p.activeConns, conn.clientAddr)
					} else {
						// 放回池中
						p.idleConns <- conn
						goto cleanup_done
					}
				default:
					goto cleanup_done
				}
			}
		cleanup_done:
			p.mu.Unlock()
		}
	}
}

func main() {

}
