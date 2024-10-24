package response

import (
	"context"
	"net/http"

	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xhttp"
)

func WriteError(ctx context.Context, w http.ResponseWriter, code int, err error) {
	Write(ctx, w, code, NewRESTErrorResponse(ctx, err))
}

func Write(ctx context.Context, w http.ResponseWriter, code int, resp any) {
	xcontext.SessionManager(ctx).Save(w, xcontext.Session(ctx))
	if err := xhttp.WriteResponseJSON(w, code, resp); err != nil {
		xcontext.Logger(ctx).Critical("failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, url string, code int) {
	xcontext.SessionManager(ctx).Save(w, xcontext.Session(ctx))
	http.Redirect(w, r, url, code)
}
