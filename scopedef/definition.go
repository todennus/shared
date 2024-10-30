package scopedef

import (
	"github.com/todennus/x/scope"
)

type User any
type Admin any
type App any

type DescriptiveScope interface {
	scope.Scoper
	Description() string
}

var StandardScopes []*standardScope
var UserScopes []*titledScope[User]
var AppScopes []*titledScope[App]
var AdminScopes []*titledScope[Admin]

var OrderedScopes []DescriptiveScope
var AllScopes = make(map[string]DescriptiveScope)
var ReadOnlyScopes = make(map[string]DescriptiveScope)

var Engine = scope.NewEngine()

var (
	OfflineAccess = standard("offline_access", "Maintain access to resource even if user is not present").readonly()
)

var (
	// User
	UserReadUserProfile = user("read:user.profile", "Grant read-only access to user profile").readonly()

	// Client
	UserReadClientProfile = user("read:client.profile", "Grant read-only access to client profile").readonly()
	UserCreateClient      = user("create:client", "Grant permission to create a client")
)

var (
	// Client
	AppReadClientOwner   = app("read:client.owner", "Grant read-only access to client owner id").readonly()
	AppReadClientProfile = app("read:client.profile", "Grant read-only access to client profile").readonly()
)

var (
	// User
	AdminReadUserProfile = admin("read:user.profile", "Grant admin read-only access to user profile").readonly()
	AdminValidateUser    = admin("validate:user", "Grant admin permission to validate user credentials").readonly()
	AdminCreateUser      = admin("create:user", "Grant admin permission to create a new user")

	// Client
	AdminReadClientProfile = admin("read:client.profile", "Grant admin read-only access to client profile").readonly()
	AdminValidateClient    = admin("validate:client", "Grant admin permission to validate client").readonly()
	AdminCreateClient      = admin("create:client", "Grant admin permission to create a client")
)

func user(value, description string) *titledScope[User] {
	s := defineTitledScope[User]("", value, description)
	UserScopes = append(UserScopes, s)
	return s
}

func app(value, description string) *titledScope[App] {
	s := defineTitledScope[App]("app", value, description)
	AppScopes = append(AppScopes, s)
	return s
}

func admin(value, description string) *titledScope[Admin] {
	s := defineTitledScope[Admin]("admin", value, description)
	AdminScopes = append(AdminScopes, s)
	return s
}

func defineTitledScope[T any](title, value, description string) *titledScope[T] {
	scope := scope.Define(Engine, newTitledScope[T](title, value, description))
	return define(scope)
}

func standard(value, description string) *standardScope {
	scope := scope.Define(Engine, newStandardScope(value, description))
	StandardScopes = append(StandardScopes, scope)
	return define(scope)
}

func define[S DescriptiveScope](s S) S {
	AllScopes[s.Scope()] = s
	OrderedScopes = append(OrderedScopes, s)
	return s
}
