package probe

import (
	"encoding/json"
	"net/http"

	"github.com/kamva/hexa/hlog"
)

func jsonDocsHandler(l *[]*HandlerDescriptor) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(*l)
		if err != nil {
			hlog.Error("probe document handler can not return list of probe handlers", hlog.ErrStack(err))
			_, _ = w.Write([]byte("occurred an internal error!"))
		}
	}
}
