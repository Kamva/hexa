package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamva/hexa"
)

type ServiceA struct {
}

func (s *ServiceA) Boot() error {
	fmt.Println("boot service a")
	return nil
}

type ServiceB struct {
	A *ServiceA
}

func (s *ServiceB) Boot() error {
	fmt.Println("boot service b")
	return nil
}

func (s *ServiceB) Shutdown(_ context.Context) error {
	fmt.Println("shutdown service b")
	return nil
}

type ServiceC struct {
}

func (s *ServiceC) Shutdown(_ context.Context) error {
	fmt.Println("shutdown service c")
	return errors.New("I can not shutdown service c")
}

type ServiceD struct {
}

func (s *ServiceD) Shutdown(_ context.Context) error {
	fmt.Println("shutdown service d")
	return nil
}

var _ hexa.Bootable = &ServiceA{}
var _ hexa.Bootable = &ServiceB{}
var _ hexa.Shutdownable = &ServiceB{}
var _ hexa.Shutdownable = &ServiceC{}
var _ hexa.Shutdownable = &ServiceD{}
