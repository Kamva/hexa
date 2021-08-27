package main

import (
	"fmt"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/sr"
)

func main() {
	// Please note all services do not contain boot or shutdown, so
	// you will not see boot or shutdown log for all services.

	a := &ServiceA{}
	b := &ServiceB{
		A: a,
	}
	c := &ServiceC{}
	d := &ServiceD{}

	r := sr.New()
	r.Register("a", a)
	r.Register("b", b)
	r.Register("c", c)

	gutil.PanicErr(r.Boot())

	// Register service D after boot:
	r.Register("d", d)

	fmt.Println("after service registration and boot, press ctrl+c to shutdown services")

	//go r.Shutdown(context.Background())
	go sr.ShutdownBySignals(r, time.Second*30)
	<-r.ShutdownCh()

	fmt.Println("by by :)")
}
