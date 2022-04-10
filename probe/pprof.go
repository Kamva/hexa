package probe

import (
	"net/http/pprof"
)

func RegisterPprofHandlers(s Server) {
	s.Register("pprof-index", "/debug/pprof/", pprof.Index, "pprof index")
	s.Register("pprof-cmdline", "/debug/pprof/cmdline", pprof.Cmdline, "pprof cmdline")
	s.Register("pprof-profile", "/debug/pprof/profile", pprof.Profile, "pprof profile")
	s.Register("pprof-symbol", "/debug/pprof/symbol", pprof.Symbol, "pprof symbol")
	s.Register("pprof-trace", "/debug/pprof/trace", pprof.Trace, "pprof trace")
}
