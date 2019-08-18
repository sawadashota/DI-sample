package notification

import "context"

type Repository interface {
	List(ctx context.Context, offset, limit int) ([]Notification, error)
	Add(ctx context.Context, n *Notification) error
	Get(ctx context.Context, id string) (*Notification, error)
	Update(ctx context.Context, n *Notification) error
	Delete(ctx context.Context, id string) error
}
