# mysql lock

- [12.14 Locking Functions](https://docs.oracle.com/cd/E17952_01/mysql-5.7-en/locking-functions.html)

`GET_LOCT(name,timeout)` - 定义一个名称为name,持续时间为time秒的锁.

- `name` - 锁的名称，一个字符串值，它必须唯一地标识你想要获取的锁。
- `timeout` - 等待锁的超时时间，单位是秒(s)。如果省略该参数或设置为NULL或设置为0，则表示无限等待.

返回值:
锁定成功，返回1；
如果尝试超时，返回0；
如果遇到错误，返回NULL.

`RELEASE_LOCK(name)` - 释放名称为name的锁.

-- name：一个字符串参数，指定了要释放的锁的名称。这个名称必须与你之前使用 GET_LOCK() 函数获取锁时所使用的名称完全相同。
-- 返回值：
-- 如果锁被成功释放（即，当前会话持有该锁），则返回 1。
-- 如果锁不存在（即，当前会话没有持有该锁），则返回 0。
-- 如果发生错误（如连接断开）等解锁失败，则返回 NULL。

`IS_FREE_LOCK(name)` - 判断是否已使用了名称为name的锁.`

如果已使用，返回0；
如果未使用，返回1；

## 示例

假设有两个会话（会话A和会话B）想要操作同一个共享资源，并且需要确保在同一时间只有一个会话可以操作该资源。

会话A:

```sql
SET @lock_result = GET_LOCK('my_lock', 10);  
IF @lock_result = 1 THEN  
    -- 锁已获取，安全地执行你的操作  
    SELECT 'Lock acquired, performing critical section' AS Status;  
    -- 完成后释放锁  
    DO RELEASE_LOCK('my_lock');  
ELSE  
    -- 无法获取锁，可能因为其他会话已持有该锁  
    SELECT 'Unable to acquire lock' AS Status;  
END IF;

```

会话B:

几乎与会话A的示例相同，但如果会话A已经持有了锁，则会话B将在GET_LOCK()调用中等待最多10秒（或者直到会话A释放锁）。

注意事项

- 锁是跨会话的，但不跨MySQL服务器实例。
- 如果MySQL服务器重启，所有锁都将被释放。
- GET_LOCK() 提供的锁是一种咨询锁（advisory lock），它不会阻止其他用户读取或写入锁定的表。它仅仅用于在应用程序逻辑中实现同步。
- 在高并发的环境下，锁可能会成为性能瓶颈。因此，应当谨慎使用，并考虑是否有更高效的同步机制可用。
- GET_LOCK() 和 RELEASE_LOCK() 的使用应当在应用程序层面仔细设计，以确保它们被正确使用且不会导致死锁或其他并发问题。

## Dev & Test

通过 docker compose 启动 mysql 数据库

```
docker compose up
# or
docker compose up -d
```

运行单元测试

```
go test . -v
# or
go test . -v -run TestLock
```
