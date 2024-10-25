package middleware

import (
	"fmt"
	"net/http"

	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/response"
)

func Recoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
						response.RESTWriteError(
							r.Context(),
							w,
							http.StatusInternalServerError,
							errordef.ErrServer.Hide(fmt.Errorf("%v", rvr), "unexpected-panic"),
						)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
