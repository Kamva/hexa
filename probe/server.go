package probe

import (
	"context"
	"net/http"
	"sync"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type Handler http.HandlerFunc

type HandlerDescriptor struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Description string  `json:"description"`
	Handler     Handler `json:"-"`
}

type Server interface {
	hexa.Runnable
	hexa.Shutdownable
	// Register registers request handler. if needed we can add
	// support to get middlewares... as options too, but currently
	// we don't need to it.
	Register(name, path string, handler Handler, description string)
}

type probeServer struct {
	server      *http.Server
	srvMux      *http.ServeMux
	mux         sync.Mutex
	descriptors []*HandlerDescriptor
}

func NewProbeServer(server *http.Server, mux *http.ServeMux) Server {
	server.Handler = mux
	pserver := &probeServer{
		server: server,
		srvMux: mux,
	}

	// Register probe server docs handler
	pserver.Register("docs", "/", jsonDocsHandler(&pserver.descriptors), "show probe server documents")

	return pserver
}

func (s *probeServer) Run() error {
	go func() {
		hlog.Info("start the probe server", hexa.StringField("address", s.server.Addr))
		err := s.server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				hlog.Info("The probe server is closed")
				return
			}

			hlog.Error("error on health check server", hexa.ErrStackField(tracer.Trace(err)), hexa.ErrField(err))
		}
	}()

	return nil
}

func (s *probeServer) Register(name, pattern string, handler Handler, description string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.descriptors = append(s.descriptors, &HandlerDescriptor{
		Name:        name,
		Path:        pattern,
		Description: description,
		Handler:     handler,
	})

	s.srvMux.HandleFunc(pattern, handler)
}

func (s *probeServer) Shutdown(c context.Context) error {
	return s.server.Shutdown(c)
}

var _ Server = &probeServer{}
