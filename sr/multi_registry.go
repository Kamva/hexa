package sr

import "github.com/kamva/hexa"

// multiSearchRegistry is an implementation of registry which can search in multiple registries
// when we want to fetch a service.
type multiSearchRegistry struct {
	hexa.ServiceRegistry
	search []hexa.ServiceRegistry
}

func (m *multiSearchRegistry) Descriptors() []*hexa.Descriptor {
	l := make([]*hexa.Descriptor, 0)
	for i := len(m.search) - 1; i > 0; i-- { // get descriptor list from registries in revert order.
		l = append(l, m.search[i].Descriptors()...)
	}

	return l
}

func (m *multiSearchRegistry) Descriptor(name string) *hexa.Descriptor {
	for _, r := range m.search {
		if d := r.Descriptor(name); d != nil {
			return d
		}
	}

	return nil
}

func (m *multiSearchRegistry) Service(name string) hexa.Service {
	if d := m.Descriptor(name); d != nil {
		return d.Instance
	}
	return nil
}

// NewMultiSearchRegistry returns new instance of multiSearchRegistry which can search in multiple registries.
func NewMultiSearchRegistry(primary hexa.ServiceRegistry, search ...hexa.ServiceRegistry) hexa.ServiceRegistry {
	return &multiSearchRegistry{
		ServiceRegistry: primary,
		search:          search,
	}
}

var _ hexa.ServiceRegistry = &multiSearchRegistry{}
