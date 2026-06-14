package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/robin38n/reqviz/backend/internal/proxy"
)

// Client-facing messages for proxy failures. They never expose the underlying
// reason (resolved IPs, DNS internals, SSRF detection); the real error is logged.
const (
	msgHostNotAllowed  = "This host isn't in the spec's approved list. You can only call the hosts you approved for this spec."
	msgRateLimited     = "Too many requests to this host. Wait a few seconds and try again."
	msgSSRFBlocked     = "That address can't be reached. ReqViz only allows requests to public internet hosts, not local or private networks."
	msgProxyBadRequest = "The request couldn't be built. Check the URL, method, and body."
	msgProxyInternal   = "Something went wrong while sending the request."
)

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ValidationError{Error: &msg})
}

// writeProxyError maps a proxy execution error to its HTTP status and safe
// client message, logging the real cause for blocked and internal failures.
func writeProxyError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, proxy.ErrHostNotAllowed):
		writeError(w, http.StatusForbidden, msgHostNotAllowed)
	case errors.Is(err, proxy.ErrRateLimited):
		writeError(w, http.StatusTooManyRequests, msgRateLimited)
	case errors.Is(err, proxy.ErrSSRFBlocked):
		slog.Default().Warn("proxy: request blocked", "error", err)
		writeError(w, http.StatusForbidden, msgSSRFBlocked)
	case errors.Is(err, proxy.ErrBadRequest):
		slog.Default().Warn("proxy: bad request", "error", err)
		writeError(w, http.StatusBadRequest, msgProxyBadRequest)
	default:
		slog.Default().Error("proxy: internal error", "error", err)
		writeError(w, http.StatusInternalServerError, msgProxyInternal)
	}
}
