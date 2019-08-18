package driver

type Driver interface {
	Registry() Registry
	CallRegistry() Driver
}

type DefaultDriver struct {
	r Registry
}

func NewDefaultDriver() Driver {
	return &DefaultDriver{
		r: NewDefaultRegistry(),
	}
}

func (d *DefaultDriver) Registry() Registry {
	return d.r
}

func (d *DefaultDriver) CallRegistry() Driver {
	callRegistry(d.r)
	return d
}

func callRegistry(r Registry) {
	r.AdminRepository()
	r.HealthHandler()
	r.JSONWebKeySetRepository()
	r.AuthHandler()
	r.NotificationRepository()
	r.NotificationHandler()
}
