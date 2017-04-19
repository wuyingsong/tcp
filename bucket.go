package tcp

import (
	"sync"
)

type TcpConnBucket struct {
	m  map[string]*TcpConn
	mu *sync.RWMutex
}

func newTcpConnBucket() *TcpConnBucket {
	return &TcpConnBucket{
		m:  make(map[string]*TcpConn),
		mu: new(sync.RWMutex),
	}
}

func (b *TcpConnBucket) Put(key string, c *TcpConn) {
	b.mu.Lock()
	b.m[key] = c
	b.mu.Unlock()
}

func (b *TcpConnBucket) Get(key string) *TcpConn {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if conn, ok := b.m[key]; ok {
		return conn
	}
	return nil
}

func (b *TcpConnBucket) Delete(key string) {
	b.mu.Lock()
	delete(b.m, key)
	b.mu.Unlock()
}
func (b *TcpConnBucket) GetAll() map[string]*TcpConn {
	b.mu.RLock()
	defer b.mu.RUnlock()
	m := make(map[string]*TcpConn, len(b.m))
	for k, v := range b.m {
		m[k] = v
	}
	return m
}
