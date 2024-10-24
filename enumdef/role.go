package enumdef

import "github.com/todennus/x/enum"

type UserRole int
type UserRoleEnum enum.Enum[UserRole]

var (
	UserRoleAdmin = enum.New[UserRole](1, "admin")
	UserRoleUser  = enum.New[UserRole](2, "user")
)
