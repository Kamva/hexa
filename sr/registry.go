// Package sr implements service registry.
// This package does not provide dependency injection,
// for that you need to use a package like google wire.
package sr

import (
	"context"
	"reflect"
	"sort"
	"sync/atomic"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

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

type serviceRegistry struct {
	l []*Descriptor

	booted     uint32 // is 1 if you boot services.
	done       uint32 // is 1 if you shutdown services.
	shutdownCh chan struct{}
}

func New() ServiceRegistry {
	return &serviceRegistry{
		l:          make([]*Descriptor, 0),
		shutdownCh: make(chan struct{}),
	}
}

func (r *serviceRegistry) Register(name string, instance Service) {
	r.register(&Descriptor{
		Name:     name,
		Instance: instance,
		Priority: len(r.l),
	})
}

func (r *serviceRegistry) RegisterByInstance(instance Service) {
	r.register(&Descriptor{
		Name:     reflect.TypeOf(instance).Elem().Name(),
		Instance: instance,
		Priority: len(r.l),
	})
}

func (r *serviceRegistry) register(d *Descriptor) {
	if r.Service(d.Name) != nil {
		hlog.Warn("you are overwriting service in service registry", hlog.String("name", d.Name))
	}

	if _, ok := d.Instance.(Bootable); ok && atomic.LoadUint32(&r.booted) == 1 {
		hlog.Debug("new registered service is bootable, but you booted your services, so the new service will not boot automatically",
			hlog.String("name", d.Name))
	}

	r.l = append(r.l, d)
}

func (r *serviceRegistry) Boot() error {
	if !atomic.CompareAndSwapUint32(&r.booted, 0, 1) {
		hlog.Warn("skip service registry boot, it has been ran already!")
		return nil
	}

	for _, d := range r.Descriptors() {
		bootable, ok := d.Instance.(Bootable)
		if !ok {
			continue
		}
		log := hlog.With(hlog.String("name", d.Name), hlog.Int("priority", d.Priority))

		log.Debug("boot service")
		if err := bootable.Boot(); err != nil {
			log.Error("service boot failed")
			return tracer.Trace(err)
		}
	}

	return nil
}

func (r *serviceRegistry) Shutdown(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&r.done, 0, 1) {
		go func() {
			dl := r.Descriptors()
			// sort descending.
			sort.Slice(dl, func(i int, j int) bool { return dl[i].Priority > dl[j].Priority })
			for _, d := range dl {
				shutdownable, ok := d.Instance.(Shutdownable)
				if !ok {
					continue
				}

				log := hlog.With(hlog.String("name", d.Name), hlog.Int("priority", d.Priority))
				log.Debug("shutdown service")
				if err := shutdownable.Shutdown(ctx); err != nil {
					log.Error("failed service shutdown")
				}
			}

			close(r.shutdownCh)
		}()
	} else {
		hlog.Debug("skip service registry shutdown and just wait to shutdown services, it has been ran already!")
	}

	select {
	case <-ctx.Done():
		hlog.Error(`shutdown context timed out, we can not shutdown remained services`)
		return tracer.Trace(ctx.Err())
	case <-r.shutdownCh:
		hlog.Info("app shutdown.")
		return nil
	}
}

func (r *serviceRegistry) ShutdownCh() (shutdownCh chan struct{}) {
	return r.shutdownCh
}

func (r *serviceRegistry) Descriptors() []*Descriptor {
	sort.Slice(r.l, func(i, j int) bool { return r.l[i].Priority < r.l[j].Priority })
	return r.l
}

func (r *serviceRegistry) Descriptor(name string) *Descriptor {
	for _, d := range r.l {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (r *serviceRegistry) Service(name string) Service {
	if d := r.Descriptor(name); d != nil {
		return d.Instance
	}
	return nil
}

var _ ServiceRegistry = &serviceRegistry{}
