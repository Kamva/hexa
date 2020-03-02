package mgmadapter

import (
	"errors"
	"github.com/Kamva/kitty"
	"github.com/Kamva/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// kittyID is the implementation of kitty ID interface for the mgm package.
type kittyID struct {
	id primitive.ObjectID
}

func (i *kittyID) String() string {
	return i.id.Hex()
}

func (i *kittyID) From(val interface{}) error {
	if oid, ok := val.(primitive.ObjectID); ok {
		i.id = oid
		return nil
	}

	if idStr, ok := val.(string); ok {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return tracer.Trace(errors.New("id value is invalid and con not covert it to primitive.ObjectID"))
		}

		i.id = id
		return nil
	}

	return tracer.Trace(errors.New("id type is invalid and can not convert it to primitive.ObjectID"))
}

func (i *kittyID) MustFrom(id interface{}) {
	if err := i.From(id); err != nil {
		panic(err)
	}
}

func (i *kittyID) Val() interface{} {
	return i.id
}

// IDD function get an id and returns IDD
func KittyID(id interface{}) kitty.ID {
	i := &kittyID{}
	i.MustFrom(id)
	return i
}

var _ kitty.ID = &kittyID{}
