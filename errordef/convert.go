package errordef

import (
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/todennus/x/xerror"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func ConvertGormError(err error) error {
	switch {
	case xerror.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDuplicated
	default:
		return err
	}
}

func ConvertRedisError(err error) error {
	switch {
	case xerror.Is(err, redis.Nil):
		return ErrNotFound
	default:
		return err
	}
}

func ConvertGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	msg := st.Message()
	code, description, found := strings.Cut(msg, ":")
	if !found {
		return err
	}

	switch code {
	case "server_error":
		return xerror.Enrich(ErrServer, description)
	case "server_timeout":
		return xerror.Enrich(ErrServerTimeout, description)
	case "invalid_request":
		return xerror.Enrich(ErrRequestInvalid, description)
	case "duplicated":
		return xerror.Enrich(ErrDuplicated, description)
	case "not_found":
		return xerror.Enrich(ErrNotFound, description)
	case "invalid_credentials":
		return xerror.Enrich(ErrCredentialsInvalid, description)
	case "unauthenticated":
		return xerror.Enrich(ErrUnauthenticated, description)
	case "forbidden":
		return xerror.Enrich(ErrForbidden, description)
	case "invalid_client":
		return xerror.Enrich(ErrClientInvalid, description)
	case "invalid_scope":
		return xerror.Enrich(ErrScopeInvalid, description)
	case "access_denined":
		return xerror.Enrich(ErrAccessDenied, description)
	case "invalid_grant":
		return xerror.Enrich(ErrTokenInvalidGrant, description)
	default:
		return err
	}
}
