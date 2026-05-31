package server

import (
	"sync"

	"github.com/hashicorp/yamux"
)

type Registry struct {
	mu      sync.RWMutex
	clients map[string]*yamux.Session
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[string]*yamux.Session),
	}
}

func (r *Registry) Add(subhost string, session *yamux.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[subhost] = session
}

func (r *Registry) Get(subhost string) (*yamux.Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.clients[subhost]
	return s, ok
}

func (r *Registry) Remove(subhost string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, subhost)
}
