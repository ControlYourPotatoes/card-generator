package memory

import (
	"fmt"
	"sync"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
	"github.com/ControlYourPotatoes/card-generator/internal/store"
)

// MemoryStore implements store.Store interface with in-memory storage
type MemoryStore struct {
	cards map[string]card.Card
	mutex sync.RWMutex
}

// New creates a new memory-based store
func New() store.Store {
	return &MemoryStore{
		cards: make(map[string]card.Card),
	}
}

// generateID creates a unique identifier for a card
func generateID(c card.Card) string {
	return fmt.Sprintf("%s-%s", c.GetType(), c.GetName())
}

func (s *MemoryStore) Save(c card.Card) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate card before saving
	if err := c.Validate(); err != nil {
		return "", fmt.Errorf("invalid card: %w", err)
	}

	id := generateID(c)
	s.cards[id] = c
	return id, nil
}

func (s *MemoryStore) Load(id string) (card.Card, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	c, exists := s.cards[id]
	if !exists {
		return nil, fmt.Errorf("card not found: %s", id)
	}
	return c, nil
}

func (s *MemoryStore) List() ([]card.Card, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	cards := make([]card.Card, 0, len(s.cards))
	for _, c := range s.cards {
		cards = append(cards, c)
	}
	return cards, nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.cards[id]; !exists {
		return fmt.Errorf("card not found: %s", id)
	}
	delete(s.cards, id)
	return nil
}

func (s *MemoryStore) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cards = nil
	return nil
}