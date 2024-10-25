package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
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

					logEntry := middleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rvr, debug.Stack())
					} else {
						middleware.PrintPrettyStack(rvr)
					}

					if r.Header.Get("Connection") != "Upgrade" {
						response.RESTWriteError(r.Context(), w, http.StatusInternalServerError, errordef.ErrServer)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
