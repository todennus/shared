package middleware

import (
	"mime"
	"net/http"
	"time"

	"github.com/todennus/shared/config"
	"github.com/todennus/shared/xcontext"
)

func LogRequest(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			originContentType := r.Header.Get("Content-Type")
			contentType, _, err := mime.ParseMediaType(originContentType)
			if err != nil {
				contentType = originContentType
			}

			if contentType == "" {
				contentType = "<empty>"
			}

			xcontext.Logger(ctx).Debug(
				"request",
				"uri", r.RequestURI,
				"method", r.Method,
				"content_type", contentType,
				"rip", r.RemoteAddr,
				"node_id", config.Variable.Server.NodeID,
			)

			start := time.Now()
			next.ServeHTTP(w, r)

			xcontext.Logger(ctx).Debug("response", "rtt", time.Since(start))
		})
	}
}
