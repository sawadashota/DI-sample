package driver

import (
	"github.com/sawadashota/di-sample/auth"
	"github.com/sawadashota/di-sample/internal/admin"
	"github.com/sawadashota/di-sample/notification"
)

type RegistryMemory struct {
	*RegistryBase
}

func NewRegistryMemory() *RegistryMemory {
	r := &RegistryMemory{
		RegistryBase: new(RegistryBase),
	}

	r.RegistryBase.with(r)

	return r
}

func (r *RegistryMemory) AdminRepository() admin.Repository {
	if r.RegistryBase.adr == nil {
		r.adr = admin.NewMemoryRepository()
	}
	return r.RegistryBase.adr
}

func (r *RegistryMemory) NotificationRepository() notification.Repository {
	if r.RegistryBase.nor == nil {
		r.nor = notification.NewMemoryRepository()
	}
	return r.RegistryBase.nor
}

func (r *RegistryMemory) JSONWebKeySetRepository() auth.Repository {
	if r.RegistryBase.aur == nil {
		r.aur = auth.NewMemoryRepository()
	}
	return r.RegistryBase.aur
}
