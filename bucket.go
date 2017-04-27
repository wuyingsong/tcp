package tcp

import (
	"sync"
	"time"
)

//TCPConnBucket 用来存放和管理TCPConn连接
type TCPConnBucket struct {
	m  map[string]*TCPConn
	mu *sync.RWMutex
}

func NewTCPConnBucket() *TCPConnBucket {
	tcb := &TCPConnBucket{
		m:  make(map[string]*TCPConn),
		mu: new(sync.RWMutex),
	}
	tcb.removeClosedTCPConnLoop()
	return tcb
}

func (b *TCPConnBucket) Put(id string, c *TCPConn) {
	b.mu.Lock()
	b.m[id] = c
	b.mu.Unlock()
}

func (b *TCPConnBucket) Get(id string) *TCPConn {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if conn, ok := b.m[id]; ok {
		return conn
	}
	return nil
}

func (b *TCPConnBucket) Delete(id string) {
	b.mu.Lock()
	delete(b.m, id)
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

func (b *TCPConnBucket) removeClosedTCPConnLoop() {
	go func() {
		removeKey := make(map[string]struct{})
		for key, conn := range b.GetAll() {
			if conn.IsClosed() {
				removeKey[key] = struct{}{}
			}
		}
		for key := range removeKey {
			b.Delete(key)
		}
		time.Sleep(time.Millisecond * 100)
	}()
}
