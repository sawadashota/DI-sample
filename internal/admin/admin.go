package admin

import (
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	password []byte
}

func (a *Admin) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword(a.password, []byte(password))
}

func (a *Admin) UpdatePassword(raw string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	a.password = hashed
	return nil
}
