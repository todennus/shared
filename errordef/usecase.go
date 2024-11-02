package errordef

import (
	"errors"
	"fmt"

	"github.com/todennus/x/xerror"
)

var (
	ErrServer        = xerror.Enrich(unique.New("server_error"), "an unexpected error occurred")
	ErrServerTimeout = xerror.Enrich(unique.New("server_timeout"), "server timeout")

	ErrRequestInvalid  = unique.New("invalid_request")
	ErrRequestTooLarge = unique.New("too_large_request")
	ErrDuplicated      = unique.New("duplicated")
	ErrNotFound        = unique.New("not_found")

	ErrCredentialsInvalid = unique.New("invalid_credentials")

	ErrUnauthenticated = unique.New("unauthenticated")
	ErrForbidden       = unique.New("forbidden")

	ErrClientInvalidType = unique.New("invalid_client_type")

	// File error
	ErrFileMismatchedSize = unique.New("mismatched_file_size")
	ErrFileMismatchedType = unique.New("mismatched_file_type")
	ErrFileInvalidContent = unique.New("invalid_file_content")

	// OAuth2 flow error
	ErrOAuth2ClientInvalid = unique.New("invalid_client")
	ErrOAuth2ScopeInvalid  = unique.New("invalid_scope")
	ErrOAuth2AccessDenied  = unique.New("access_denied")
	ErrOAuth2InvalidGrant  = unique.New("invalid_grant")
)

// For handling domain error
var DomainWrapper = xerror.NewWrapperConfigs(ErrServer, ErrDomainKnown)

var unique = uniqueError{}

type uniqueError map[string]error

func (u uniqueError) New(s string) error {
	if _, ok := u[s]; ok {
		panic(fmt.Sprintf("duplicated error %s", s))
	}

	u[s] = errors.New(s)
	return u[s]
}
