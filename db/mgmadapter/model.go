//go:generate msgp

package mgmadapter

// Entity struct contain model's default fields.
type Entity struct {
	IDField `bson:",inline"`
}

// TouchableEntity struct is the entity base struct with touch feature
// to increase its version on each change or set date fields on each
// update.
type TouchableEntity struct {
	IDField   `bson:",inline"`
	Touchable `bson:",inline"`
}