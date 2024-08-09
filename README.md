# lock

定义一个锁接口 `Locker`，通过不同组件来实现该锁，以便用户业务实际演进情况。

> 核心目标不是性能。而是基于用户当前业务场景已有的组件，实现一个分布式锁，以便服务演进。

## lock tables

- redis
- mysql
- etcd
- consul
- zookeeper
- kafka
- ...
