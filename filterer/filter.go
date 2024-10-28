package filterer

import (
	"context"
	"reflect"

	"github.com/todennus/shared/scopedef"
	"github.com/todennus/x/scope"
	"github.com/todennus/x/xcontext"
	"github.com/xybor-x/snowflake"
)

type Filterer[T any] struct {
	ctx      context.Context
	obj      *T
	filtered bool
	old      T
	target   T
}

func Set[T any](ctx context.Context, obj *T, target T) *Filterer[T] {
	return &Filterer[T]{ctx: ctx, obj: obj, old: *obj, target: target}
}

func Filter[T any](ctx context.Context, obj *T) *Filterer[T] {
	var t T
	return Set(ctx, obj, t)
}

func (f *Filterer[T]) When(b func(ctx context.Context) bool) *Filterer[T] {
	if !f.filtered && b(f.ctx) {
		f.set()
	}
	return f
}

func (f *Filterer[T]) WhenNot(b func(ctx context.Context) bool) *Filterer[T] {
	return f.When(func(ctx context.Context) bool { return !b(ctx) })
}

func (f *Filterer[T]) WhenNotContainsScope(target scope.Scope) *Filterer[T] {
	return f.WhenNot(func(ctx context.Context) bool { return xcontext.Scope(ctx).Contains(target) })
}

func (f *Filterer[T]) WhenNotContainsScopeWithUserID(target scope.Scope, userID snowflake.ID) *Filterer[T] {
	if !xcontext.Scope(f.ctx).Contains(target.WithTitle(scopedef.TitleAdmin)) {
		return f.When(func(ctx context.Context) bool { return xcontext.RequestUserID(ctx) != userID })
	}

	return f.WhenNotContainsScope(target)
}

func (f *Filterer[T]) set() {
	reflect.ValueOf(f.obj).Elem().Set(reflect.ValueOf(f.target))
}
