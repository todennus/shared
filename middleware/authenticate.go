package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/todennus/shared/response"
	"github.com/todennus/shared/scopedef"
	"github.com/todennus/shared/tokendef"
	"github.com/todennus/x/token"
	"github.com/todennus/x/xcontext"
)

func WithAuthenticate(ctx context.Context, authorization string, engine token.Engine) context.Context {
	if authorization == "" {
		return ctx
	}

	tokenType, token, found := strings.Cut(authorization, " ")
	if !found {
		return ctx
	}

	if engine.Type() != tokenType {
		return ctx
	}

	accessToken := tokendef.OAuth2AccessToken{}
	ok, err := engine.Validate(ctx, token, &accessToken)
	if err != nil {
		xcontext.Logger(ctx).Debug("failed-to-parse-token", "err", err)
		return ctx
	}

	if !ok {
		xcontext.Logger(ctx).Debug("expired token")
		return ctx
	}

	ctx = xcontext.WithRequestUserID(ctx, accessToken.SnowflakeSub())
	ctx = xcontext.WithScope(ctx, scopedef.Engine.ParseScopes(accessToken.Scope))

	xcontext.Logger(ctx).Debug("auth-info", "uid", accessToken.Subject, "scope", accessToken.Scope)

	return ctx
}

func Authentication(engine token.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authorization := r.Header.Get("Authorization")

			next.ServeHTTP(w, r.WithContext(WithAuthenticate(ctx, authorization, engine)))
		})
	}
}

func RequireAuthentication(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if xcontext.RequestUserID(ctx) == 0 {
			response.Write(ctx, w, http.StatusUnauthorized,
				response.NewRESTErrorResponseWithMessage(ctx, "unauthenticated", "require authentication to access api"))
		} else {
			handler(w, r)
		}
	}
}
