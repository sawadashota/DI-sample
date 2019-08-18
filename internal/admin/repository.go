package admin

import "context"

type Registry interface {
	AdminRepository() Repository
}

type Repository interface {
	List(ctx context.Context) ([]Admin, error)
	Find(ctx context.Context, id string) (*Admin, error)
	FindByEmail(ctx context.Context, email string) (*Admin, error)
	Add(ctx context.Context, admin *Admin) error
	Cancel(ctx context.Context, id string) error
}
