package kitty

type User interface {
	IsGuest() bool
	GetID() interface{}
}
