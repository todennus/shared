package enumdef

import "github.com/todennus/x/enum"

type UserRole int

var (
	UserRoleAdmin = enum.New[UserRole](1, "admin")
	UserRoleUser  = enum.New[UserRole](2, "user")
)
