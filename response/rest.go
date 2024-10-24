package response

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/todennus/shared/errordef"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xerror"
	"github.com/todennus/x/xhttp"
)

const TimeLayout = "2024-10-20T15:45:30Z"

type RESTResponseStatus string

const (
	RESTResponseStatusSuccess RESTResponseStatus = "success"
	RESTResponseStatusError                      = "error"
)

type RESTMetadata struct {
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}

func NewRESTMetadata(ctx context.Context) *RESTMetadata {
	return &RESTMetadata{
		Timestamp: time.Now(),
		RequestID: xcontext.RequestID(ctx),
	}
}

type RESTResponse struct {
	Status           RESTResponseStatus `json:"status,omitempty"`
	Data             any                `json:"data,omitempty"`
	Error            string             `json:"error,omitempty"`
	ErrorDescription string             `json:"error_description,omitempty"`
	Metadata         *RESTMetadata      `json:"metadata,omitempty"`
}

func NewRESTResponse(data any) *RESTResponse {
	return &RESTResponse{
		Status: RESTResponseStatusSuccess,
		Data:   data,
	}
}

func NewRESTErrorResponse(ctx context.Context, err error) *RESTResponse {
	var richError xerror.RichError
	if errors.As(err, &richError) {
		if richError.Detail() != nil {
			attrs := []any{"err", richError.Detail()}
			attrs = append(attrs, richError.Attributes()...)
			if errors.Is(err, errordef.ErrServer) {
				xcontext.Logger(ctx).Warn(richError.Event(), attrs...)
			} else {
				xcontext.Logger(ctx).Debug(richError.Event(), attrs...)
			}
		}

		return NewRESTErrorResponseWithMessage(ctx, richError.Code().Error(), richError.Description())
	}

	xcontext.Logger(ctx).Critical("internal-error", "err", err)
	return NewRESTUnexpectedErrorResponse(ctx)
}

func NewRESTUnexpectedErrorResponse(ctx context.Context) *RESTResponse {
	return NewRESTErrorResponseWithMessage(ctx,
		"server_error",
		"an unexpected error occurred, please contact to admin if you see this error",
	)
}

func NewRESTErrorResponseWithMessage(ctx context.Context, err string, description string) *RESTResponse {
	response := &RESTResponse{Status: RESTResponseStatusError, Metadata: NewRESTMetadata(ctx)}

	response.Error = err
	response.ErrorDescription = description

	return response
}

func SetQuery(ctx context.Context, q url.Values, err error) {
	errResp := NewRESTErrorResponse(ctx, err)

	q.Set("error", errResp.Error)
	q.Set("error_description", errResp.ErrorDescription)
	q.Set("timestamp", errResp.Metadata.Timestamp.Format(TimeLayout))
	q.Set("request_id", errResp.Metadata.RequestID)
}

func RESTWriteLogInvalidRequestError(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		panic("do not pass a nil error here")
	}

	var code int
	response := &RESTResponse{}
	switch {
	case xerror.Is(err, xhttp.ErrHTTPBadRequest, errordef.ErrRequestInvalid):
		code = http.StatusBadRequest
		response = NewRESTErrorResponseWithMessage(ctx, "invalid_request", err.Error())
	default:
		code = http.StatusInternalServerError
		response = NewRESTUnexpectedErrorResponse(ctx)

		xcontext.Logger(ctx).Debug("failed-to-parse-data", "err", err)
	}

	Write(ctx, w, code, response)
}

type RESTResponseHandler struct {
	err         error
	resp        any
	code        int
	defaultCode int
}

func NewRESTResponseHandler(ctx context.Context, resp any, err error) *RESTResponseHandler {
	if timeoutErr := context.Cause(ctx); timeoutErr != nil && errors.Is(timeoutErr, errordef.ErrServerTimeout) {
		err = errordef.ErrServerTimeout.Hide(err, "timeout")
	}

	return (&RESTResponseHandler{
		err:  err,
		resp: resp,
		code: -1,
	}).WithDefaultCode(http.StatusOK).Map(http.StatusGatewayTimeout, errordef.ErrServerTimeout)
}

func (h *RESTResponseHandler) WithDefaultCode(code int) *RESTResponseHandler {
	h.defaultCode = code
	return h
}

func (h *RESTResponseHandler) Map(code int, errs ...error) *RESTResponseHandler {
	if h.err == nil || h.code != -1 {
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

func (h *RESTResponseHandler) WriteHTTPResponse(ctx context.Context, w http.ResponseWriter) {
	h.Map(http.StatusInternalServerError)

	if h.code == -1 {
		h.code = h.defaultCode
	}

	var resp any
	if h.err != nil {
		resp = NewRESTErrorResponse(ctx, h.err)
	} else {
		resp = NewRESTResponse(h.resp)
	}

	Write(ctx, w, h.code, resp)
}

func (h *RESTResponseHandler) WriteHTTPResponseWithoutWrap(ctx context.Context, w http.ResponseWriter) {
	h.Map(http.StatusInternalServerError)

	if h.code == -1 {
		h.code = h.defaultCode
	}

	var resp any
	if h.err != nil {
		errResp := NewRESTErrorResponse(ctx, h.err)
		errResp.Status = ""
		resp = errResp
	} else {
		resp = h.resp
	}

	Write(ctx, w, h.code, resp)
}

func (h *RESTResponseHandler) Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, code int) {
	h.Map(http.StatusInternalServerError)

	if h.code == -1 {
		h.code = code
	}

	if h.err != nil {
		Write(ctx, w, h.code, NewRESTErrorResponse(ctx, h.err))
		return
	}

	Redirect(ctx, w, r, h.resp.(string), h.code)
}
