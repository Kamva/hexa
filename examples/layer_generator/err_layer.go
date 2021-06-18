package main

type errLayer struct {
	app App
}

func (a *errLayer) By() (string, error) {
	r1, r2 := a.app.By()

	if r2 != nil {
		// do something when error occured.
	}

	return r1, r2

}

func (a *errLayer) ByFromB() (string, error) {
	r1, r2 := a.app.ByFromB()

	if r2 != nil {
		// do something when error occured.
	}

	// Do something, this method has "trace" annotation

	return r1, r2

}

func (a *errLayer) Hi(username string, password string) (sentence string, err error) {
	r1, r2 := a.app.Hi(username, password)

	if r2 != nil {
		// do something when error occured.
	}

	return r1, r2

}

func NewErrLayer(app App) App {
	return &errLayer{
		app: app,
	}
}
