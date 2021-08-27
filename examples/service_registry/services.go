package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamva/hexa/sr"
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

var _ sr.Bootable = &ServiceA{}
var _ sr.Bootable = &ServiceB{}
var _ sr.Shutdownable = &ServiceB{}
var _ sr.Shutdownable = &ServiceC{}
var _ sr.Shutdownable = &ServiceD{}
