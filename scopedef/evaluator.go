package scopedef

import (
	"context"
	"reflect"

	"github.com/todennus/shared/xcontext"
	"github.com/todennus/x/scope"
	"github.com/xybor-x/snowflake"
)

type evaluator struct {
	scopes      scope.Scopes
	isSatisfied bool
}

func Eval(scopes scope.Scopes) *evaluator {
	return &evaluator{scopes: scopes, isSatisfied: false}
}

type readyEvaluator struct {
	*evaluator
}

func (e *evaluator) RequireAdmin(scope *titledScope[Admin]) *readyEvaluator {
	if e.isSatisfied {
		return &readyEvaluator{e}
	}

	if e.scopes.Contains(scope) {
		e.isSatisfied = true
	}

	return &readyEvaluator{e}
}

func (e *evaluator) RequireAnyUser(scope *titledScope[User]) *readyEvaluator {
	if e.isSatisfied {
		return &readyEvaluator{e}
	}

	if e.scopes.Contains(scope) {
		e.isSatisfied = true
	}

	return &readyEvaluator{e}
}

func (e *evaluator) RequireUser(ctx context.Context, scope *titledScope[User], userID snowflake.ID) *readyEvaluator {
	if e.isSatisfied {
		return &readyEvaluator{e}
	}

	if e.scopes.Contains(scope) && xcontext.RequestSubjectID(ctx) == userID {
		e.isSatisfied = true
	}

	return &readyEvaluator{e}
}

func (e *evaluator) RequireAnyApp(scope *titledScope[App]) *readyEvaluator {
	if e.isSatisfied {
		return &readyEvaluator{e}
	}

	if e.scopes.Contains(scope) {
		e.isSatisfied = true
	}

	return &readyEvaluator{e}
}

func (e *evaluator) RequireApp(ctx context.Context, scope *titledScope[App], clientID snowflake.ID) *readyEvaluator {
	if e.isSatisfied {
		return &readyEvaluator{e}
	}

	if e.scopes.Contains(scope) && xcontext.RequestSubjectID(ctx) == clientID {
		e.isSatisfied = true
	}

	return &readyEvaluator{e}
}

func (e *readyEvaluator) IsSatisfied() bool {
	return e.isSatisfied
}

func (e *readyEvaluator) IsUnsatisfied() bool {
	return !e.IsSatisfied()
}

func (e *readyEvaluator) SetIfUnsatisfied(obj, value any) {
	if !e.isSatisfied {
		reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(value))
	}
}

func (a *readyEvaluator) FilterIfUnsatisfied(obj ...any) {
	if !a.isSatisfied {
		for i := range obj {
			reflect.ValueOf(obj[i]).Elem().Set(reflect.Zero(reflect.TypeOf(obj[i]).Elem()))
		}
	}
}
