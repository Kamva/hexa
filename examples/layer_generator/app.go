package main

import "fmt"

type App interface {
	Hi(username string, password string) (sentence string, err error)
	By() (string, error)
}

type appImpl struct {
}

func (a *appImpl) Hi(username string, password string) (sentence string, err error) {
	return fmt.Sprintf("hi %s", username), nil
}

func (a *appImpl) By() (string, error) {
	return "By", nil
}

var _ App = &appImpl{}
