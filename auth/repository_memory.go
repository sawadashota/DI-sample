package auth

import (
	"context"
	"log"
	"sync"

	"github.com/pkg/errors"
)

type MemoryRepository struct {
	jsonWebKeySets []JSONWebKeySet
	mu             sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	r := &MemoryRepository{
		jsonWebKeySets: make([]JSONWebKeySet, 0),
		mu:             sync.RWMutex{},
	}

	_ = r.seed()

	return r
}

func (r *MemoryRepository) seed() error {
	for i := 0; i < 3; i++ {
		set, err := NewJSONWebKeySet()
		if err != nil {
			return err
		}
		if err := r.Add(context.Background(), set); err != nil {
			return err
		}
	}
	return nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]JSONWebKeySet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.jsonWebKeySets, nil
}

func (r *MemoryRepository) First(ctx context.Context) (*JSONWebKeySet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.jsonWebKeySets) == 0 {
		return nil, errors.New("no json web key set data")
	}

	set := r.jsonWebKeySets[0]
	return &set, nil
}

func (r *MemoryRepository) Find(ctx context.Context, kid string) (*JSONWebKeySet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, set := range r.jsonWebKeySets {
		log.Println(set.ID)
		if set.ID == kid {
			return &set, nil
		}
	}

	return nil, errors.Errorf("key id %s is not found", kid)
}

func (r *MemoryRepository) Add(ctx context.Context, key *JSONWebKeySet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.jsonWebKeySets = append([]JSONWebKeySet{*key}, r.jsonWebKeySets...)
	return nil
}

func (r *MemoryRepository) Delete(ctx context.Context, kid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []JSONWebKeySet
	for _, set := range r.jsonWebKeySets {
		if set.ID != kid {
			result = append(result, set)
		}
	}

	if len(result) == len(r.jsonWebKeySets) {
		return errors.Errorf("key id %s is not found", kid)
	}
	r.jsonWebKeySets = result

	return nil
}
