package hexa

import "context"

type Service any // Currently Service interface does not needs to implement anything.

type Bootable interface {
	Boot() error
}

// Runnable is for services that need to be ran in background.
type Runnable interface {
	// Run MUST run in non-blocking mode.
	// The first return param is a channel which you can
	// close it to specify the run is done. or send an error
	// to it to specify you have error in your run (and then you can close it).
	// The second param is an error param which specifies you didn't run the app
	// and had an error before running it.
	//
	//
	// - How to check if the service is running? when the service returns from this method, it
	//   specifies that the service is running.
	// - How to specify service had error before start running? the second return param is the
	//   error which specifies that.
	// - How to specify if the service had an error when it was running? the service will send
	//   an error in the close channel and then closes it.
	// - How to specify if the service run is done? it closes the close channel.
	//
	// Previous signature of the `Run` method was `Run() error`, and behaviour(blocking or non-blocking) was unspecified.
	// In some of our runnable services we needed a way to know if a service is running before start running our
	// integration tests. For example, we needed to make sure the HTTP server is started running before sending request
	// to it in our tests, but we didn't have any way to know it, so we specified this signature to let services return
	// from the `Run()` method when they ran the service and also return a channel to specify if the
	// Run is done (e.g., server is sopped). and also in this way all services are just in one mode: non-blocking.
	Run() (close <-chan error, err error)
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
