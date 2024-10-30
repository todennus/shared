package scopedef

import (
	"github.com/todennus/x/scope"
)

const namespace = "todennus/"

type titledScope[T any] struct {
	description string

	title string
	inner *scope.Scope
}

func newTitledScope[T any](title, value, description string) *titledScope[T] {
	return &titledScope[T]{
		title:       title,
		description: description,
		inner:       scope.New(value),
	}
}

func (s *titledScope[T]) Scope() string {
	result := namespace

	if s.title != "" {
		result += s.title + ":"
	}

	return result + s.inner.Scope()
}

func (s *titledScope[T]) Description() string {
	return s.description
}

func (s *titledScope[T]) readonly() *titledScope[T] {
	ReadOnlyScopes[s.Scope()] = s
	return s
}

type standardScope struct {
	inner       *scope.Scope
	description string
}

func newStandardScope(value, description string) *standardScope {
	return &standardScope{
		description: description,
		inner:       scope.New(value),
	}
}

func (s *standardScope) Scope() string {
	return s.inner.Scope()
}

func (s *standardScope) Description() string {
	return s.description
}

func (s *standardScope) readonly() *standardScope {
	ReadOnlyScopes[s.Scope()] = s
	return s
}
