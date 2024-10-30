package authentication

import (
	"context"
	"strings"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/shared/scopedef"
	"github.com/todennus/shared/tokendef"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/x/token"
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
	if err := engine.Validate(ctx, token, &accessToken); err != nil {
		xcontext.Logger(ctx).Debug("failed-to-parse-token", "err", err)
		return ctx
	}

	subType := ""
	if accessToken.Role != "" {
		subType = "user"
		ctx = xcontext.WithRequestSubjectType(ctx, enumdef.SubjectUser)
	} else {
		subType = "client"
		ctx = xcontext.WithRequestSubjectType(ctx, enumdef.SubjectClient)
	}

	ctx = xcontext.WithRequestSubjectID(ctx, accessToken.SnowflakeSub())
	ctx = xcontext.WithScope(ctx, scopedef.Engine.ParseAnyScopes(accessToken.Scope))

	xcontext.Logger(ctx).Debug("auth-info",
		"sub", accessToken.Subject, "scope", accessToken.Scope, "type", subType)

	return ctx
}
