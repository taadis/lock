package lock

import (
	"context"
	"log/slog"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlLock struct {
	key        string
	expiration time.Duration
	db         *sql.DB
}

func NewMysqlLock(db *sql.DB, key string, expiration time.Duration) *MysqlLock {
	return &MysqlLock{
		key:        key,
		expiration: expiration,
		db:         db,
	}
}

func (l *MysqlLock) Lock(ctx context.Context) (bool, error) {
	result, err := l.getLock(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "GET_LOCK() failed", "key", l.key, "error", err)
		return false, err
	}

	slog.InfoContext(ctx, "GET_LOCK", "name", l.key, "result.Valid", result.Valid, "result.Int64", result.Int64)

	// 1表示获取锁成功
	if result.Valid && result.Int64 == 1 {
		return true, nil
	}

	// 0表示获取锁失败(因为超时?)
	return false, nil
}

func (l *MysqlLock) Unlock(ctx context.Context) (bool, error) {
	result, err := l.releaseLock(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "RELEASE_LOCK() failed", "key", l.key, "error", err)
		return false, err
	}

	slog.InfoContext(ctx, "RELEASE_LOCK()", "key", l.key, "result.Valid", result.Valid, "result.Int64", result.Int64)

	if !result.Valid {
		slog.WarnContext(ctx, "RELEASE_LOCK() failed", "key", l.key, "error", "result is NULL")
		return false, nil
	}
	// 1表示释放锁成功
	if result.Valid && result.Int64 == 1 {
		return true, nil
	}
	// 0表示释放锁失败(锁不存在?)
	if result.Valid && result.Int64 == 0 {
		slog.WarnContext(ctx, "RELEASE_LOCK() failed", "key", l.key, "error", "lock not exist")
		return false, nil
	}

	return false, nil
}

func (l *MysqlLock) WaitLock(ctx context.Context) error {
	for {
		result, err := l.isFreeLock(ctx)
		if err != nil {
			return err
		}

		if result.Valid && result.Int64 == 1 {
			locked, err := l.Lock(ctx)
			if err != nil {
				return err
			}
			if locked {
				return nil
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// getLock 从 MySQL 数据库获取锁
//
// 参数：
// ctx：上下文对象，用于控制函数执行过程中的超时、取消等
//
// 返回值：
// sql.NullInt64：MySQL GET_LOCK() 函数的返回值，如果获取成功返回 1，如果获取失败返回 0
// 1表示获取锁成功
// 0表示获取锁失败(因为超时?)
// error：如果函数执行过程中发生错误，则返回相应的错误信息，否则返回 nil
func (l *MysqlLock) getLock(ctx context.Context) (sql.NullInt64, error) {
	// GET_LOCK(name,timeout)
	query := `select GET_LOCK(?, NULL)`
	var result sql.NullInt64
	err := l.db.QueryRowContext(ctx, query, l.key).Scan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (l *MysqlLock) releaseLock(ctx context.Context) (sql.NullInt64, error) {
	query := `select RELEASE_LOCK(?)`
	var result sql.NullInt64
	err := l.db.QueryRowContext(ctx, query, l.key).Scan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (l *MysqlLock) isFreeLock(ctx context.Context) (sql.NullInt64, error) {
	query := `select IS_FREE_LOCK(?)`
	var result sql.NullInt64
	err := l.db.QueryRowContext(ctx, query, l.key).Scan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (l *MysqlLock) isUsedLock(ctx context.Context) (sql.NullInt64, error) {
	query := `select IS_USED_LOCK(?)`
	var result sql.NullInt64
	err := l.db.QueryRowContext(ctx, query, l.key).Scan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
