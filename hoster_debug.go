package gohost

import (
	"errors"
	"net/http"

	// register debug http handlers
	_ "expvar"
	_ "net/http/pprof"
)

// serveDebug will start the debug endpoint.
func (h *Hoster) serveDebug() error {
	// check if debug endpoint is enabled
	if h.DebugAddr != "" {
		// validate parameters
		if h.DebugAddr == "" {
			return errors.New("debug address cannot be empty")
		}

		return http.ListenAndServe(h.DebugAddr, nil)
	}

	return nil
}
