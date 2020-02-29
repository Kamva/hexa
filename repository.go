package kitty

// Repository interface must implement by each repo.
type Repository interface {
	ValidateID(val interface{}) error
	PrepareID(val interface{}) (ID interface{}, err error)
	MustPrepareID(val interface{}) (ID interface{})
}

