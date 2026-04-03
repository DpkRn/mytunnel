package server

import (
	"sync"

	"github.com/hashicorp/yamux"
)

// Registry maps tunnel subdomain ids to yamux sessions.
type Registry struct {
	mu      sync.RWMutex
	clients map[string]*yamux.Session
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[string]*yamux.Session),
	}
}

func (r *Registry) Add(id string, s *yamux.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[id] = s
}

func (r *Registry) Get(id string) (*yamux.Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.clients[id]
	return s, ok
}

func (r *Registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, id)
}
