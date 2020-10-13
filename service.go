// service.go
package greeting

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/rs/xid"
)

// GreetingService for Greetings
type GreetingService interface {
	GetAllForUser(ctx context.Context, username string) ([]Greeting, error)
	GetByID(ctx context.Context, id string) (Greeting, error)
	Add(ctx context.Context, greeting Greeting) (Greeting, error)
	Update(ctx context.Context, id string, greeting Greeting) error
	Delete(ctx context.Context, id string) error
}

// *** Implementation ***

var (
	// ErrIncosistentIDs is when the ID given in payload differs from params ID
	ErrInconsistentIDs = errors.New("Incosistent IDs")
	// ErrNotFound is when the Entity doesn't exist
	ErrNotFound = errors.New("Not found")
)

// NewInmemGreetingService creates an in memory Greeting service
func NewInmemGreetingService() GreetingService {
	s := &inmemService{
		m: map[string]Greeting{},
	}
	rand.Seed(time.Now().UnixNano())
	return s
}

// inmemService is an In Memor
type inmemService struct {
	sync.RWMutex
	m map[string]Greeting
}

// GetAllForUser gets Greetings from memory for a user
func (s *inmemService) GetAllForUser(ctx context.Context, username string) ([]Greeting, error) {
	s.RLock()
	defer s.RUnlock()

	greetings := make([]Greeting, 0, len(s.m))
	for _, greeting := range s.m {
		if greeting.Username == username {
			greetings = append(greetings, greeting)
		}
	}
	return greetings, nil
}

// Get an Greetings from the database
func (s *inmemService) GetByID(ctx context.Context, id string) (Greeting, error) {
	s.Lock()
	defer s.Unlock()

	if greeting, ok := s.m[id]; ok {
		return greeting, nil
	}

	return Greeting{}, nil
}

// Add a Greeting to memory
func (s *inmemService) Add(ctx context.Context, greeting Greeting) (Greeting, error) {
	s.Lock()
	defer s.Unlock()

	greeting.ID = xid.New().String()
	greeting.CreatedOn = time.Now()

	s.m[greeting.ID] = greeting
	return greeting, nil
}

// Update a Greeting in memory
func (s *inmemService) Update(ctx context.Context, id string, greeting Greeting) error {
	s.Lock()
	defer s.Unlock()

	if id != greeting.ID {
		return ErrInconsistentIDs
	}

	s.m[greeting.ID] = greeting
	return nil
}

// Delete a Greeting from memory
func (s *inmemService) Delete(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}

	delete(s.m, id)
	return nil
}
