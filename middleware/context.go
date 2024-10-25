package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/todennus/shared/config"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xcrypto"
)

var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func WithBasicContext(ctx context.Context, config *config.Config) context.Context {
	ctx = xcontext.WithRequestID(ctx, xcrypto.RandString(16))
	ctx = xcontext.WithSessionManager(ctx, config.SessionManager)
	ctx = xcontext.WithLogger(ctx, config.Logger.With("request_id", xcontext.RequestID(ctx)))

	return ctx
}

func SetupContext(config *config.Config) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rip := realIP(r); rip != "" {
				r.RemoteAddr = rip
			}

			h.ServeHTTP(w, r.WithContext(WithBasicContext(r.Context(), config)))
		})
	}
}

func realIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
