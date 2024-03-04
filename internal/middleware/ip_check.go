package middleware

import (
	"net"
	"net/http"

	"github.com/0loff/go_link_shortener/internal/logger"
	"go.uber.org/zap"
)

// IPChecker is the closure function for init middleware with external params
func IPChecker(trusted string) func(h http.Handler) http.Handler {
	IPC := NewIPCheck(trusted)
	return IPC.Handler
}

// IPCheck is a new middleware structure for ip checking
type IPCheck struct {
	trusted string
}

// NewIPCheck - initialization constructor for new middleware
func NewIPCheck(t string) *IPCheck {
	return &IPCheck{trusted: t}
}

// Handler - main middleware method
func (ipc *IPCheck) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ipc.trusted == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, trustedSubnet, err := net.ParseCIDR(ipc.trusted)
		if err != nil {
			logger.Log.Error("The value of the trusted subnet could not be parsed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		requestIP := net.ParseIP(r.Header.Get("X-Real-IP"))

		if !trustedSubnet.Contains(requestIP) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
