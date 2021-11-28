package main

import (
	"fmt"
	"net/http"
	"syscall"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/probe"
)

const addr = "localhost:7676"

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	probeServer := probe.NewServer(server, mux)

	probeServer.Register("hi", "/hi", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hi"))
		if err != nil {
			panic(err)
		}
	}, "hi server")

	gutil.PanicErr(probeServer.Run())
	fmt.Println("server is listening on", addr)
	gutil.WaitForSignals(syscall.SIGINT)
}
