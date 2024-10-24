package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/todennus/config"
	"github.com/todennus/shared/errordef"
)

func Timeout(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithoutCancel(ctx)

			timeout := time.Duration(config.Variable.Server.RequestTimeout) * time.Millisecond
			ctx, cancel := context.WithTimeoutCause(ctx, timeout, errordef.ErrServerTimeout)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
