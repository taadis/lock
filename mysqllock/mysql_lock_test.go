package lock

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestMysqlLock(t *testing.T) {
	dsn := "dev:dev@tcp(127.0.0.1:3306)/dev"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open mysql error:%+v", err)
	}
	defer db.Close()

	// 验证链接
	err = db.Ping()
	if err != nil {
		t.Fatalf("ping mysql error:%+v", err)
	}
	t.Log("Successfully connected to mysql!")

	ctx := context.Background()

	t.Run("incorrect_lock_name", func(t *testing.T) {
		l := NewMysqlLock(db, "", time.Second)
		result, err := l.Lock(ctx)
		if err != nil {
			if err.Error() != "Error 3057 (42000): Incorrect user-level lock name ''." {
				t.Fatalf("lock error:%+v", err)
			}
		}

		if result {
			t.Fatalf("invalid lock result:%v", result)
		}
	})

	t.Run("lock_success", func(t *testing.T) {
		l := NewMysqlLock(db, "test:mysql:lock", time.Second)
		result, err := l.Lock(ctx)
		if err != nil {
			t.Fatalf("lock error:%+v", err)
		}
		if !result {
			t.Fatalf("invalid lock result:%v", result)
		}
	})

	t.Run("unlock", func(t *testing.T) {
		l := NewMysqlLock(db, "test:mysql:lock", time.Second)
		result, err := l.Unlock(ctx)
		if err != nil {
			t.Fatalf("unlock error:%+v", err)
		}
		if !result {
			t.Fatalf("invalid unlock result:%v", result)
		}
		// todo:这里为什么一直释放不掉锁/直接执行sql需要返回NULL时才行?
	})

	t.Run("wait lock", func(t *testing.T) {
		key := "test:mysql:lock"
		l := NewMysqlLock(db, key, time.Second)
		err := l.WaitLock(ctx)
		if err != nil {
			t.Fatalf("lock erro1r:%+v", err)
		}
		time.Sleep(10 * time.Second)
		result, err := l.Unlock(ctx)
		if err != nil {
			t.Fatalf("unlock error:%+v", err)
		}
		if !result {
			t.Fatalf("invalid unlock result:%v", result)
		}
	})
}
