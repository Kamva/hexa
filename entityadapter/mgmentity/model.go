package mgmentity

// KittyModel struct contain model's default fields.
type KittyModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

// Creating function call to it's inner fields defined hooks
func (model *KittyModel) Creating() error {
	return model.DateFields.Creating()
}

// Saving function call to it's inner fields defined hooks
func (model *KittyModel) Saving() error {
	return model.DateFields.Saving()
}

