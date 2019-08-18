package auth

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/google/uuid"
)

type JSONWebKeySet struct {
	ID      string
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func NewJSONWebKeySet() (*JSONWebKeySet, error) {
	const keyLength = 4096

	key, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, err
	}

	return &JSONWebKeySet{
		ID:      uuid.New().String(),
		Private: key,
		Public:  &key.PublicKey,
	}, nil
}
