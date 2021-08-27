package hexa

import "context"

type Service interface{} // Currently Service interface does not needs to implement anything.

type Bootable interface {
	Boot() error
}

// Runnable is for services that need to be ran in background.
type Runnable interface {
	Run() error
}

type Shutdownable interface {
	Shutdown(context.Context) error
}

// Descriptor describes the service.
type Descriptor struct {
	Name     string
	Instance Service
	Priority int
}

type ServiceRegistry interface {
	Register(name string, instance Service)
	RegisterByInstance(instance Service)
	Boot() error
	Shutdown(ctx context.Context) error
	ShutdownCh() chan struct{}

	// Descriptors returns descriptors ordered by their priority.
	Descriptors() []*Descriptor
	Descriptor(name string) *Descriptor
	// Service method should return nil if service not found.
	Service(name string) Service
}
