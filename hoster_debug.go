package gohost

import (
	"net/http"

	// register debug http handlers
	_ "expvar"
	_ "net/http/pprof"
)

// serveDebug will start the debug endpoint.
func (h *Hoster) serveDebug() error {
	// check if debug endpoint is enabled
	if h.DebugAddr != "" {
		return http.ListenAndServe(h.DebugAddr, nil)
	}

	return nil
}
