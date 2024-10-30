package xcontext

import (
	"context"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/x/logging"
	"github.com/todennus/x/scope"
	"github.com/todennus/x/session"
	"github.com/xybor-x/snowflake"
)

type contextKey int

const (
	loggerKey contextKey = iota
	requestIDKey
	requestSubjectIDKey
	requestSubjectTypeKey
	scopeKey
	sessionKey
	sessionManagerKey
)

func WithLogger(ctx context.Context, logger logging.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) logging.Logger {
	if val := ctx.Value(loggerKey); val != nil {
		return val.(logging.Logger)
	}

	return logging.NewSLogger(logging.LevelDebug).With("logger", "temporary")
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestID(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

func WithRequestSubjectID(ctx context.Context, userID snowflake.ID) context.Context {
	return context.WithValue(ctx, requestSubjectIDKey, userID)
}

func RequestSubjectID(ctx context.Context) snowflake.ID {
	if val := ctx.Value(requestSubjectIDKey); val != nil {
		return val.(snowflake.ID)
	}

	return 0
}

func WithRequestSubjectType(ctx context.Context, t enumdef.SubjectType) context.Context {
	return context.WithValue(ctx, requestSubjectTypeKey, t)
}

func RequestSubjectType(ctx context.Context) enumdef.SubjectType {
	if val := ctx.Value(requestSubjectTypeKey); val != nil {
		return val.(enumdef.SubjectType)
	}

	return -1
}

func WithScope(ctx context.Context, scopes scope.Scopes) context.Context {
	return context.WithValue(ctx, scopeKey, scopes)
}

func Scope(ctx context.Context) scope.Scopes {
	if val := ctx.Value(scopeKey); val != nil {
		return val.(scope.Scopes)
	}

	return nil
}

func WithSession(ctx context.Context, session *session.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func Session(ctx context.Context) *session.Session {
	if val := ctx.Value(sessionKey); val != nil {
		return val.(*session.Session)
	}

	return &session.Session{}
}

func WithSessionManager(ctx context.Context, manager *session.Manager) context.Context {
	return context.WithValue(ctx, sessionManagerKey, manager)
}

func SessionManager(ctx context.Context) *session.Manager {
	return ctx.Value(sessionManagerKey).(*session.Manager)
}
