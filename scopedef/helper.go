package scopedef

import "github.com/todennus/x/scope"

func IsTitle[T any](scope scope.Scoper) bool {
	_, ok := scope.(*titledScope[T])
	return ok
}

func HasAnyTitle[T any](scope scope.Scopes) bool {
	for i := range scope {
		if IsTitle[T](scope[i]) {
			return true
		}
	}

	return false
}

func HasAnyStandard(scope scope.Scopes) bool {
	for i := range scope {
		if _, ok := scope[i].(*standardScope); ok {
			return true
		}
	}

	return false
}

func HasAnyExternal(scopes scope.Scopes) bool {
	for i := range scopes {
		if _, ok := AllScopes[scopes[i].Scope()]; !ok {
			return true
		}
	}

	return false
}

func IsAllReadonly(scopes scope.Scopes) bool {
	for i := range scopes {
		if _, ok := ReadOnlyScopes[scopes[i].Scope()]; !ok {
			return false
		}
	}

	return true
}
