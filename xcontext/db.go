package xcontext

import (
	"context"

	"gorm.io/gorm"
)

type DBTransaction struct {
	db *gorm.DB
}

func WithDBTransaction(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbtxKey, &DBTransaction{})
}

func DB(ctx context.Context, df *gorm.DB) *gorm.DB {
	if val := ctx.Value(dbtxKey); val != nil {
		dbtx := val.(*DBTransaction)
		if dbtx.db == nil {
			dbtx.db = df.Begin()
		}

		return dbtx.db.WithContext(ctx)
	}

	return df.WithContext(ctx)
}

func DBCommit(ctx context.Context) context.Context {
	if val := ctx.Value(dbtxKey); val != nil {
		dbtx := val.(*DBTransaction)
		if dbtx.db != nil {
			if err := dbtx.db.Commit().Error; err != nil {
				Logger(ctx).Warn("failed-to-commit-transaction", "err", err)
			}
		}
	}

	return context.WithValue(ctx, dbtxKey, nil)
}

func DBRollback(ctx context.Context) context.Context {
	if val := ctx.Value(dbtxKey); val != nil {
		dbtx := val.(*DBTransaction)
		if dbtx.db != nil {
			if err := dbtx.db.Rollback(); err != nil {
				Logger(ctx).Warn("failed-to-rollback-transaction", "err", err)
			}
		}
	}

	return context.WithValue(ctx, dbtxKey, nil)
}
