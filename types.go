package kitty

type Secret string

func (p Secret) String() string {
	return "****"
}


