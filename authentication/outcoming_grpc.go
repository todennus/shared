package authentication

import (
	"context"
	"fmt"

	"github.com/todennus/shared/xcontext"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"
)

type GrpcAuthorization struct {
	t      *oauth2.Token
	source func(context.Context) oauth2.TokenSource
}

func NewGrpcAuthorization(source func(context.Context) oauth2.TokenSource) *GrpcAuthorization {
	return &GrpcAuthorization{source: source}
}

func (a *GrpcAuthorization) Context(ctx context.Context) context.Context {
	if !a.t.Valid() {
		var err error
		a.t, err = oauth2.ReuseTokenSource(a.t, a.source(ctx)).Token()
		if err != nil {
			xcontext.Logger(ctx).Warn("failed-to-get-token-source", "err", err)
			return ctx
		}
	}

	return metadata.NewOutgoingContext(
		ctx,
		metadata.MD{"authorization": {fmt.Sprintf("%s %s", a.t.TokenType, a.t.AccessToken)}},
	)
}
