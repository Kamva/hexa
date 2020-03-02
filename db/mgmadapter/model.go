package mgmadapter

import "github.com/Kamva/tracer"

// Entity struct contain model's default fields.
type Entity struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

// Creating function call to it's inner fields defined hooks
func (model *Entity) Creating() error {
	return tracer.Trace(model.DateFields.Creating())
}

// Saving function call to it's inner fields defined hooks
func (model *Entity) Saving() error {
	return tracer.Trace(model.DateFields.Saving())
}
