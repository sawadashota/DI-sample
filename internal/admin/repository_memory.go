package admin

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MemoryRepository struct {
	admin []Admin
	mu    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	r := &MemoryRepository{
		admin: make([]Admin, 0),
		mu:    sync.RWMutex{},
	}

	_ = r.seed()

	return r
}

func (r *MemoryRepository) seed() error {
	admin := &Admin{
		ID:    uuid.New().String(),
		Name:  "Shota",
		Email: "example@example.com",
	}

	if err := admin.UpdatePassword("test"); err != nil {
		return nil
	}

	if err := r.Add(context.Background(), admin); err != nil {
		return nil
	}

	return nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Admin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.admin, nil
}

func (r *MemoryRepository) Find(ctx context.Context, id string) (*Admin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, admin := range r.admin {
		if admin.ID == id {
			return &admin, nil
		}
	}

	return nil, errors.Errorf("admin id %s is not found", id)
}

func (r *MemoryRepository) FindByEmail(ctx context.Context, email string) (*Admin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, admin := range r.admin {
		if admin.Email == email {
			return &admin, nil
		}
	}

	return nil, errors.Errorf("admin email %s is not found", email)
}

func (r *MemoryRepository) Add(ctx context.Context, admin *Admin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.admin = append(r.admin, *admin)
	return nil
}

func (r *MemoryRepository) Cancel(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []Admin
	for _, admin := range r.admin {
		if admin.ID != id {
			result = append(result, admin)
		}
	}

	if len(result) == len(r.admin) {
		return errors.Errorf("admin id %s is not found", id)
	}
	r.admin = result

	return nil
}
