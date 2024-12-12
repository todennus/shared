package enumdef

import (
	"github.com/xybor-x/enum"
)

type userRole any
type UserRole = enum.WrapEnum[userRole]

const (
	UserRoleUser UserRole = iota
	UserRoleAdmin
)

func init() {
	enum.Map(UserRoleUser, "user")
	enum.Map(UserRoleAdmin, "admin")
	enum.Finalize[UserRole]()
}
