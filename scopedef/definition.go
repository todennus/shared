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
	OfflineAccess = standard("offline_access", "Maintain access to the resource even if the user is not present").readonly()
)

var (
	UserReadUserProfile  = user("read:user.profile", "Grant read-only access to the user's profile").readonly()
	UserReadUserAvatar   = user("read:user.avatar", "Grant read-only access to the user's avatar").readonly()
	UserUpdateUserAvatar = user("update:user.avatar", "Grant permission to update the user's avatar")

	UserReadClientProfile = user("read:client.profile", "Grant read-only access to the client's profile").readonly()
	UserCreateClient      = user("create:client", "Grant permission to create new clients")
)

var (
	AppReadClientOwner   = app("read:client.owner", "Grant read-only access to the client's owner id").readonly()
	AppReadClientProfile = app("read:client.profile", "Grant read-only access to the client's profile").readonly()
)

var (
	AdminReadUserProfile = admin("read:user.profile", "Grant read-only access to all users' profiles").readonly()
	AdminValidateUser    = admin("validate:user", "Grant permission to validate all users' credentials")
	AdminCreateUser      = admin("create:user", "Grant permission to create new users")

	AdminReadClientProfile = admin("read:client.profile", "Grant read-only access to all client profiles").readonly()
	AdminValidateClient    = admin("validate:client", "Grant permission to validate the client's credentials")
	AdminCreateClient      = admin("create:client", "Grant permission to create new clients")

	AdminRegisterFilePolicy          = admin("register:file.policy", "Grant permission to register file upload policy")
	AdminCreatePresignedFile         = admin("create:file.presigned_url", "Grant permission to create file presigned url")
	AdminChangeRefcountFileOwnership = admin("change:file.refcount", "Grant permission to change ref count of file ownership")
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
