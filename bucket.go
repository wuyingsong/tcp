package tcp

import (
	"sync"
)

type TCPConnBucket struct {
	m  map[string]*TCPConn
	mu *sync.RWMutex
}

func newTCPConnBucket() *TCPConnBucket {
	return &TCPConnBucket{
		m:  make(map[string]*TCPConn),
		mu: new(sync.RWMutex),
	}
}

func (b *TCPConnBucket) Put(key string, c *TCPConn) {
	b.mu.Lock()
	b.m[key] = c
	b.mu.Unlock()
}

func (b *TCPConnBucket) Get(key string) *TCPConn {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if conn, ok := b.m[key]; ok {
		return conn
	}
	return nil
}

func (b *TCPConnBucket) Delete(key string) {
	b.mu.Lock()
	delete(b.m, key)
	b.mu.Unlock()
}
func (b *TCPConnBucket) GetAll() map[string]*TCPConn {
	b.mu.RLock()
	defer b.mu.RUnlock()
	m := make(map[string]*TCPConn, len(b.m))
	for k, v := range b.m {
		m[k] = v
	}
	return m
}
