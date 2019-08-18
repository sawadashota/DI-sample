package notification

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MemoryRepository struct {
	notifications []Notification
	mu            sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	r := &MemoryRepository{
		notifications: make([]Notification, 0),
		mu:            sync.RWMutex{},
	}

	_ = r.seed()

	return r
}

func (r *MemoryRepository) seed() error {
	notifications := []Notification{
		{
			ID:        uuid.New().String(),
			Title:     "Hello Admin Console",
			Body:      "Admin Console is useful",
			IsDraft:   false,
			PublishAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Title:     "Hello Admin Console",
			Body:      "Admin Console is useful",
			IsDraft:   true,
			PublishAt: time.Now(),
		},
	}

	for _, notification := range notifications {
		if err := r.Add(context.Background(), &notification); err != nil {
			return err
		}
	}

	return nil
}

func (r *MemoryRepository) List(ctx context.Context, offset, limit int) ([]Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Notification
	for i, notification := range r.notifications {
		if i < offset {
			continue
		}
		if i >= offset+limit {
			break
		}
		result = append(result, notification)
	}

	return result, nil
}

func (r *MemoryRepository) Add(ctx context.Context, n *Notification) error {
	if _, err := r.Get(ctx, n.ID); err == nil {
		return errors.Errorf("notifications id %s is already exist", n.ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.notifications = append(r.notifications, *n)
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, notification := range r.notifications {
		if notification.ID == id {
			return &notification, nil
		}
	}

	return nil, errors.Errorf("notifications id %s is not found", id)
}

func (r *MemoryRepository) Update(ctx context.Context, n *Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, notification := range r.notifications {
		if notification.ID == n.ID {
			r.notifications[i] = *n
			return nil
		}
	}
	return errors.Errorf("notifications id %s is not found", n.ID)
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []Notification
	for _, notification := range r.notifications {
		if notification.ID != id {
			result = append(result, notification)
		}
	}

	if len(result) == len(r.notifications) {
		return errors.Errorf("notifications id %s is not found", id)
	}
	r.notifications = result

	return nil
}
