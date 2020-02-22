package kitty

type Password string

func (p Password) String() string {
	return "****"
}


