package interceptor

import (
	"context"
	"time"

	"github.com/todennus/shared/config"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/middleware"
	"github.com/todennus/x/token"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xcrypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UnaryInterceptor struct {
	basicContext bool
	timeout      bool
	authenticate bool
	logrtt       bool
}

func NewUnaryInterceptor() *UnaryInterceptor {
	return &UnaryInterceptor{}
}

func (i *UnaryInterceptor) WithBasicContext() *UnaryInterceptor {
	i.basicContext = true
	return i
}

func (i *UnaryInterceptor) WithTimeout() *UnaryInterceptor {
	i.timeout = true
	return i
}

func (i *UnaryInterceptor) WithAuthenticate() *UnaryInterceptor {
	i.authenticate = true
	return i
}

func (i *UnaryInterceptor) WithLogRoundTripTime() *UnaryInterceptor {
	i.logrtt = true
	return i
}

func (i *UnaryInterceptor) Interceptor(config *config.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if i.basicContext {
			ctx = middleware.WithBasicContext(ctx, config)
		}

		xcontext.Logger(ctx).Debug(
			"rpc_request",
			"function", info.FullMethod,
			"node_id", config.Variable.Server.NodeID,
		)

		if i.timeout {
			var cancel context.CancelFunc
			ctx, cancel = withTimeout(ctx, config)
			defer cancel()
		}

		if i.authenticate {
			ctx = withAuthenticate(ctx, config.TokenEngine)
		}

		start := time.Now()
		resp, err := handler(ctx, req)

		if i.logrtt {
			xcontext.Logger(ctx).Debug("rpc_response", "rtt", time.Since(start))
		}

		return resp, err
	}
}

func withRequestID(ctx context.Context) context.Context {
	ctx = xcontext.WithRequestID(ctx, xcrypto.RandString(16))
	logger := xcontext.Logger(ctx).With("request_id", xcontext.RequestID(ctx))
	ctx = xcontext.WithLogger(ctx, logger)
	return ctx
}

func withTimeout(ctx context.Context, config *config.Config) (context.Context, context.CancelFunc) {
	ctx = context.WithoutCancel(ctx)
	timeout := time.Duration(config.Variable.Server.RequestTimeout) * time.Millisecond
	return context.WithTimeoutCause(ctx, timeout, errordef.ErrServerTimeout)
}

func withAuthenticate(ctx context.Context, engine token.Engine) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		xcontext.Logger(ctx).Debug("not-found-metadata")
		return ctx
	}

	authorization := md["authorization"]
	if len(authorization) != 1 {
		xcontext.Logger(ctx).Debug("invalid-or-not-found-authorization-metadata")
		return ctx
	}

	return middleware.WithAuthenticate(ctx, authorization[0], engine)
}
