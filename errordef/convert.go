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

	code, description, found := strings.Cut(st.Message(), ":")
	if !found {
		return err
	}

	if e, ok := unique[code]; ok {
		return xerror.Enrich(e, description)
	}

	return err
}
