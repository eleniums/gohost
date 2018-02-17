package gohost

import (
	"net/http"
)

// serveDebug will start the debug endpoint.
func (h *Hoster) serveDebug() {
	// check if debug endpoint is enabled
	if h.DebugAddr != "" {
		go func() {
			http.ListenAndServe(h.DebugAddr, nil)
		}()
	}
}
