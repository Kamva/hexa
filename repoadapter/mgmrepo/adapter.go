package mgmrepo

import (
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/Kamva/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type adapter struct {
}

func (a adapter) ValidateID(id interface{}) error {
	if _, ok := id.(primitive.ObjectID); ok {
		return nil
	}

	if idStr, ok := id.(string); ok {
		_, err := primitive.ObjectIDFromHex(idStr)

		return tracer.Trace(err)
	}

	return tracer.Trace(errors.New("error value is invalid"))
}

func (a adapter) PrepareID(val interface{}) (ID interface{}, err error) {
	if err := a.ValidateID(val); err != nil {
		return nil, err
	}

	if idStr, ok := val.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	// Otherwise id must be ObjectId
	return val, nil
}

func (a adapter) MustPrepareID(val interface{}) (ID interface{}) {
	ID, err := a.PrepareID(val)
	gutil.PanicErr(err)
	return ID
}

var _ kitty.Repository = &adapter{}
