package errordef

import (
	"errors"
	"fmt"

	"github.com/todennus/x/xerror"
)

var (
	ErrServer        = xerror.Enrich(errors.New("server_error"), "an unexpected error occurred")
	ErrServerTimeout = xerror.Enrich(errors.New("server_timeout"), "server timeout")

	ErrRequestInvalid = errors.New("invalid_request")
	ErrDuplicated     = errors.New("duplicated")
	ErrNotFound       = errors.New("not_found")

	ErrCredentialsInvalid = errors.New("invalid_credentials")

	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")

	ErrClientInvalid = errors.New("invalid_client")

	ErrScopeInvalid = errors.New("invalid_scope")

	ErrAuthorizationAccessDenied = errors.New("access_denined")
	ErrTokenInvalidGrant         = errors.New("invalid_grant")
)

// For Domain Usage
var ErrDomainKnown = errors.New("")

func GoodDomainError(format string, a ...any) error {
	return fmt.Errorf("%w%s", ErrDomainKnown, fmt.Sprintf(format, a...))
}

func UnexpectedDomainError(format string, a ...any) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func UnexpectedDomainWrap(err error, format string, a ...any) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, a...))
}

// For handling domain error
var Domain = xerror.NewWrapperConfigs(ErrServer, ErrDomainKnown)
