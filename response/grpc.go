package response

import (
	"context"
	"errors"

	"github.com/todennus/shared/errordef"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCResponseHandler[D any] struct {
	err         error
	resp        D
	code        codes.Code
	defaultCode codes.Code
}

func NewGRPCResponseHandler[D any](ctx context.Context, resp D, err error) *GRPCResponseHandler[D] {
	if timeoutErr := context.Cause(ctx); timeoutErr != nil && errors.Is(timeoutErr, errordef.ErrServerTimeout) {
		err = errordef.ErrServerTimeout.Hide(err, "timeout")
	}

	return (&GRPCResponseHandler[D]{
		err:  err,
		resp: resp,
		code: codes.Unknown,
	}).WithDefaultCode(codes.OK).
		Map(codes.DeadlineExceeded, errordef.ErrServerTimeout)
}

func (h *GRPCResponseHandler[D]) WithDefaultCode(code codes.Code) *GRPCResponseHandler[D] {
	h.defaultCode = code
	return h
}

func (h *GRPCResponseHandler[D]) Map(code codes.Code, errs ...error) *GRPCResponseHandler[D] {
	if code == codes.Unknown {
		panic("do not use code unknown")
	}

	if h.err == nil || h.code != codes.Unknown {
		return h
	}

	if len(errs) == 0 {
		h.code = code
	} else {
		for _, err := range errs {
			if errors.Is(h.err, err) {
				h.code = code
				break
			}
		}
	}

	return h
}

func (h *GRPCResponseHandler[D]) Finalize(ctx context.Context) (D, error) {
	h.Map(codes.Internal)

	if h.code == codes.Unknown {
		h.code = h.defaultCode
	}

	var defaultResp D
	if h.err != nil {
		var richError xerror.RichError
		if errors.As(h.err, &richError) {
			if richError.Detail() != nil {
				attrs := []any{"err", richError.Detail()}
				attrs = append(attrs, richError.Attributes()...)
				if errors.Is(h.err, errordef.ErrServer) {
					xcontext.Logger(ctx).Warn(richError.Event(), attrs...)
				} else {
					xcontext.Logger(ctx).Debug(richError.Event(), attrs...)
				}
			}

			h.err = richError.Reduce()
		} else {
			xcontext.Logger(ctx).Critical("internal-error", "err", h.err)
			h.err = errors.New("unexpected_server_error: an unexpected error occured")
		}

		return defaultResp, status.Errorf(h.code, h.err.Error())
	}

	return h.resp, nil
}
