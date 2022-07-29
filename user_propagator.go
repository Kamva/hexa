package hexa

import (
	"encoding/json"

	"github.com/kamva/tracer"
)

type UserPropagator interface {
	ToBytes(User) ([]byte, error)
	FromBytes([]byte) (User, error)
}

type userPropagator struct {
}

func (p *userPropagator) ToBytes(u User) ([]byte, error) {
	return json.Marshal(u.MetaData())
}

func (p *userPropagator) FromBytes(m []byte) (User, error) {
	meta := make(map[string]any)
	if err := json.Unmarshal(m, &meta); err != nil {
		return nil, tracer.Trace(err)
	}

	if err := userMetaInterfaceToTrueTypedMeta(meta); err != nil {
		return nil, tracer.Trace(err)
	}

	return NewUserFromMeta(meta)
}

func NewUserPropagator() UserPropagator {
	return &userPropagator{}
}

var _ UserPropagator = &userPropagator{}
