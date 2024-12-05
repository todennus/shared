package tokendef

import (
	"fmt"
	"time"

	"github.com/todennus/x/token"
	"github.com/xybor-x/snowflake"
)

type FileToken struct {
	ID          string `json:"jti"`
	OwnershipID string `json:"oid"`
	FileID      string `json:"hsh"`
	UserID      string `json:"uid"`
	Type        string `json:"typ"`
	Size        int    `json:"sze"`
	ExpiresAt   int    `json:"exp"`
}

func (claims *FileToken) SnowflakeID() snowflake.ID {
	id, err := snowflake.ParseString(claims.ID)
	if err != nil {
		panic(err)
	}
	return id
}

func (claims *FileToken) SnowflakeOwnershipID() snowflake.ID {
	id, err := snowflake.ParseString(claims.OwnershipID)
	if err != nil {
		panic(err)
	}
	return id
}

func (claims *FileToken) SnowflakeUserID() snowflake.ID {
	id, err := snowflake.ParseString(claims.UserID)
	if err != nil {
		panic(err)
	}
	return id
}

func (claims *FileToken) Valid() error {
	if claims.ExpiresAt != 0 && time.Unix(int64(claims.ExpiresAt), 0).Before(time.Now()) {
		return token.ErrTokenExpired
	}

	snowflakeID, err := snowflake.ParseString(claims.ID)
	if err != nil {
		return fmt.Errorf("%w: %s", token.ErrTokenInvalidFormat, "invalid jti")
	}

	if _, err := snowflake.ParseString(claims.UserID); err != nil {
		return fmt.Errorf("%w: %s", token.ErrTokenInvalidFormat, "invalid uid")
	}

	if _, err := snowflake.ParseString(claims.OwnershipID); err != nil {
		return fmt.Errorf("%w: %s", token.ErrTokenInvalidFormat, "invalid oid")
	}

	createdAt := time.UnixMilli(snowflakeID.Time())
	if createdAt.After(time.Now()) {
		return token.ErrTokenNotYetValid
	}

	return nil
}
