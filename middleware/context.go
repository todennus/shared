package middleware

import (
	"context"
	"net/http"

	"github.com/todennus/shared/config"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xcrypto"
)

func WithBasicContext(ctx context.Context, config *config.Config) context.Context {
	ctx = xcontext.WithRequestID(ctx, xcrypto.RandString(16))
	ctx = xcontext.WithSessionManager(ctx, config.SessionManager)
	ctx = xcontext.WithLogger(ctx, config.Logger.With("request_id", xcontext.RequestID(ctx)))

	return ctx
}

func SetupContext(config *config.Config) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(WithBasicContext(r.Context(), config)))
		})
	}
}
