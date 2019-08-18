package auth

import "context"

type Repository interface {
	List(ctx context.Context) ([]JSONWebKeySet, error)
	First(ctx context.Context) (*JSONWebKeySet, error)
	Find(ctx context.Context, kid string) (*JSONWebKeySet, error)
	Add(ctx context.Context, key *JSONWebKeySet) error
	Delete(ctx context.Context, kid string) error
}
