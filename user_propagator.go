package hexa

import (
	"encoding/json"

	"github.com/kamva/gutil"
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
	meta := make(map[string]interface{})
	if err := json.Unmarshal(m, &meta); err != nil {
		return nil, tracer.Trace(err)
	}

	// Convert Usertype:
	meta[UserMetaKeyUserType] = UserType(meta[UserMetaKeyUserType].(string))

	// Convert user roles from []interface{} to []string:
	roles := make([]string, 0)
	err := gutil.UnmarshalStruct(meta[UserMetaKeyRoles], &roles)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	meta[UserMetaKeyRoles] = roles

	return NewUserFromMeta(meta)
}

func NewUserPropagator() UserPropagator {
	return &userPropagator{}
}

var _ UserPropagator = &userPropagator{}