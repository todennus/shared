package tokendef

import (
	"fmt"
	"time"

	"github.com/todennus/x/token"
	"github.com/xybor-x/snowflake"
)

var _ (token.Claims) = (*OAuth2StandardClaims)(nil)

type OAuth2StandardClaims struct {
	ID        string `json:"jti,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ExpiresAt int    `json:"exp,omitempty"`
	NotBefore int    `json:"nbf,omitempty"`
}

func (claims *OAuth2StandardClaims) SnowflakeID() snowflake.ID {
	id, err := snowflake.ParseString(claims.ID)
	if err != nil {
		panic(err)
	}
	return id
}

func (claims *OAuth2StandardClaims) SnowflakeSub() snowflake.ID {
	id, err := snowflake.ParseString(claims.Subject)
	if err != nil {
		panic(err)
	}
	return id
}

func (claims *OAuth2StandardClaims) Valid() error {
	now := time.Now()
	if claims.ExpiresAt != 0 && time.Unix(int64(claims.ExpiresAt), 0).Before(now) {
		return token.ErrTokenExpired
	}

	if claims.NotBefore != 0 && time.Unix(int64(claims.NotBefore), 0).After(now) {
		return token.ErrTokenNotYetValid
	}

	snowflakeID, err := snowflake.ParseString(claims.ID)
	if err != nil {
		return fmt.Errorf("%w: %s", token.ErrTokenInvalidFormat, "invalid jti")
	}

	if _, err := snowflake.ParseString(claims.ID); err != nil {
		return fmt.Errorf("%w: %s", token.ErrTokenInvalidFormat, "invalid sub")
	}

	createdAt := time.UnixMilli(snowflakeID.Time())
	if createdAt.After(now) {
		return token.ErrTokenNotYetValid
	}

	return nil
}

type OAuth2AccessToken struct {
	*OAuth2StandardClaims
	Scope string `json:"scope"`
}

type OAuth2RefreshToken struct {
	*OAuth2StandardClaims
	SequenceNumber int    `json:"seq"`
	Scope          string `json:"scope"`
}

type OAuth2IDToken struct {
	*OAuth2StandardClaims

	Username    string `json:"username"`
	Displayname string `json:"display_name"`
	Role        string `json:"role"`
}
