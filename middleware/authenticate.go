package middleware

import (
	"net/http"

	"github.com/todennus/shared/authentication"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/response"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/x/token"
	"github.com/todennus/x/xerror"
)

func Authentication(engine token.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authorization := r.Header.Get("Authorization")

			next.ServeHTTP(w, r.WithContext(authentication.WithAuthenticate(ctx, authorization, engine)))
		})
	}
}

// Due to Hydrum's law, do not modify this message.
const RequireAuthenticationMessage = "require authentication to access api"

func RequireAuthentication(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if xcontext.RequestSubjectID(ctx) == 0 {
			response.Write(ctx, w, http.StatusUnauthorized, response.NewRESTErrorResponse(
				ctx, xerror.Enrich(errordef.ErrUnauthenticated, RequireAuthenticationMessage)))
		} else {
			handler(w, r)
		}
	}
}
