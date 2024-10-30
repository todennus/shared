package scopedef

import (
	"errors"

	"github.com/todennus/x/scope"
)

type allower struct {
	scope         scope.Scopes
	allowUser     bool
	allowApp      bool
	allowAdmin    bool
	allowExternal bool
	allowStandard bool
}

func (a *allower) Err() error {
	if !a.allowUser && HasAnyTitle[User](a.scope) {
		return errors.New("unable to request the user scope")
	}

	if !a.allowApp && HasAnyTitle[App](a.scope) {
		return errors.New("unable to request the app scope")
	}

	if !a.allowAdmin && HasAnyTitle[Admin](a.scope) {
		return errors.New("unable to request the admin scope")
	}

	if !a.allowStandard && HasAnyStandard(a.scope) {
		return errors.New("unable to request the standard scope")
	}

	if !a.allowExternal && HasAnyExternal(a.scope) {
		return errors.New("unable to request external scope")
	}

	return nil
}

type onlyallower struct {
	*allower
}

func OnlyAllow(scopes scope.Scopes) *onlyallower {
	return &onlyallower{&allower{scope: scopes}}
}

func (a *onlyallower) HasStandard() *onlyallower {
	a.allowStandard = true
	return a
}

func (a *onlyallower) HasAdmin() *onlyallower {
	a.allowAdmin = true
	return a
}

func (a *onlyallower) HasUser() *onlyallower {
	a.allowUser = true
	return a
}

func (a *onlyallower) HasApp() *onlyallower {
	a.allowApp = true
	return a
}

func (a *onlyallower) HasExternal() *onlyallower {
	a.allowExternal = true
	return a
}

type exceptallower struct {
	*allower
}

func NotAllow(scopes scope.Scopes) *exceptallower {
	return &exceptallower{&allower{
		scope:         scopes,
		allowUser:     true,
		allowApp:      true,
		allowAdmin:    true,
		allowExternal: true,
		allowStandard: true,
	}}
}

func (a *exceptallower) HasStandard() *exceptallower {
	a.allowStandard = false
	return a
}

func (a *exceptallower) HasAdmin() *exceptallower {
	a.allowAdmin = false
	return a
}

func (a *exceptallower) HasUser() *exceptallower {
	a.allowUser = false
	return a
}

func (a *exceptallower) HasApp() *exceptallower {
	a.allowApp = false
	return a
}

func (a *exceptallower) HasExternal() *exceptallower {
	a.allowExternal = false
	return a
}
