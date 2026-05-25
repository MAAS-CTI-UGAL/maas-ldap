package users

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

var errBlankUsername = errors.New("user mapping username is required")

// Mapping stores a backend credential keyed by LDAP username.
type Mapping struct {
	Username string
	Secret   string
}

// Store keeps user mappings in memory for concurrent access.
type Store struct {
	mu       sync.RWMutex
	mappings []Mapping
}

// NewStore builds an in-memory store from existing mappings.
func NewStore(mappings []Mapping) *Store {
	store := &Store{}
	for _, mapping := range mappings {
		if err := store.Set(mapping); err != nil {
			continue
		}
	}
	return store
}

// Get returns the mapping for username.
func (s *Store) Get(username string) (Mapping, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, mapping := range s.mappings {
		if mapping.Username == username {
			return mapping, true
		}
	}
	return Mapping{}, false
}

// Set inserts or replaces one user mapping.
func (s *Store) Set(mapping Mapping) error {
	if strings.TrimSpace(mapping.Username) == "" {
		return errBlankUsername
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.mappings {
		if existing.Username == mapping.Username {
			s.mappings[i] = mapping
			return nil
		}
	}

	s.mappings = append(s.mappings, mapping)
	return nil
}

// Delete removes one user mapping.
func (s *Store) Delete(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, mapping := range s.mappings {
		if mapping.Username == username {
			s.mappings = append(s.mappings[:i], s.mappings[i+1:]...)
			return
		}
	}
}

// List returns all mappings sorted by username.
func (s *Store) List() []Mapping {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mappings := append([]Mapping(nil), s.mappings...)
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Username < mappings[j].Username
	})
	return mappings
}
