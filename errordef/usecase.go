package errordef

import (
	"errors"

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

	// OAuth2 flow error
	ErrOAuth2ClientInvalid = errors.New("invalid_client")
	ErrOAuth2ScopeInvalid  = errors.New("invalid_scope")
	ErrOAuth2AccessDenied  = errors.New("access_denied")
	ErrOAuth2InvalidGrant  = errors.New("invalid_grant")
)

// For handling domain error
var DomainWrapper = xerror.NewWrapperConfigs(ErrServer, ErrDomainKnown)
