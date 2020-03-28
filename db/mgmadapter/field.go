package mgmadapter

import (
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//IDField struct contain model's ID field.
type IDField struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

// PrepareID method prepare id value to using it as id in filtering,...
// e.g convert hex-string id value to bson.ObjectId
func (f *IDField) PrepareID(id interface{}) (objID interface{}, err error) {
	if idStr, ok := id.(string); ok {
		objID, err = primitive.ObjectIDFromHex(idStr)

		if err != nil {
			err = hexa.ErrInvalidID.SetError(tracer.Trace(err))
		}

		return objID, tracer.Trace(err)
	}

	// Otherwise id must be ObjectId
	return id, nil
}

// GetID method return model's id
func (f *IDField) GetID() interface{} {
	return f.ID
}

// SetID set id value of model's id field.
func (f *IDField) SetID(id interface{}) {
	f.ID = id.(primitive.ObjectID)
}

