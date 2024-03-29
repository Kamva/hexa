package mgmadapter

import (
	"errors"
	"fmt"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// hexaID is the implementation of hexa ID interface for the mgm package.
type hexaID struct {
	id primitive.ObjectID
}

func (i *hexaID) String() string {
	return i.id.Hex()
}

func (i *hexaID) Validate(id any) error {
	kid := &hexaID{}
	return kid.From(id)
}

func (i *hexaID) From(val any) error {
	if val == nil {
		return tracer.Trace(errors.New("id value is nil"))
	}

	if oid, ok := val.(primitive.ObjectID); ok {
		i.id = oid
		return nil
	}

	if idStr, ok := val.(string); ok {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return tracer.Trace(fmt.Errorf("id with value:%v is invalid and can not covert it to primitive.ObjectID", idStr))
		}

		i.id = id
		return nil
	}

	return tracer.Trace(errors.New("id type is invalid and can not convert it to primitive.ObjectID"))
}

func (i *hexaID) MustFrom(id any) {
	if err := i.From(id); err != nil {
		panic(err)
	}
}

func (i *hexaID) Val() any {
	return i.id
}

func (i *hexaID) IsEqual(hexaID hexa.ID) bool {
	if hexaID == nil {
		return false
	}
	if id, ok := hexaID.Val().(primitive.ObjectID); ok {
		return i.id == id
	}
	return false
}

// ID function get an id and returns IDD
func ID(id any) hexa.ID {
	i := &hexaID{}
	i.MustFrom(id)
	return i
}

// EmptyID returns empty instance of the id.
func EmptyID() hexa.ID {
	i := &hexaID{}
	return i
}

var _ hexa.ID = &hexaID{}
