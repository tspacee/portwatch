package watch

import (
	"errors"
	"fmt"
	"sync"
)

// Resolver maps port numbers to known service names.
type Resolver struct {
	mu       sync.RWMutex
	services map[int]string
}

// NewResolver returns a Resolver pre-loaded with common well-known ports.
func NewResolver() *Resolver {
	r := &Resolver{
		services: make(map[int]string),
	}
	for port, name := range wellKnownPorts {
		r.services[port] = name
	}
	return r
}

// Register associates a port with a service name.
func (r *Resolver) Register(port int, service string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %d", port)
	}
	if service == "" {
		return errors.New("service name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[port] = service
	return nil
}

// Resolve returns the service name for a port, or "unknown" if not found.
func (r *Resolver) Resolve(port int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if name, ok := r.services[port]; ok {
		return name
	}
	return "unknown"
}

// Len returns the number of registered services.
func (r *Resolver) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.services)
}

var wellKnownPorts = map[int]string{
	21:   "ftp",
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
}
