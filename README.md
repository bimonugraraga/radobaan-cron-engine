# radobaan-cron-engine
Distributed cron scheduling with Redis-backed locking. Wraps robfig/cron with seconds support and provides a lightweight lock per usecase/target so concurrent workers don’t double-execute the same job.

## Installation
- Go 1.22+
- Redis running and reachable

```bash
go get github.com/bimonugraraga/radobaan-cron-engine
```

## Quick Start
This snippet mirrors the example in [main.go](file:///Users/bytedance/go/src/github.com/radobaan-cron-engine/examples/basic/main.go).

```go
package main

import (
    "context"
    "log"
    "time"
    cronengine "github.com/bimonugraraga/radobaan-cron-engine"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()

    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatal("redis not running:", err)
    }
    scheduler := cronengine.NewScheduler(rdb)

    scheduler.Schedule(cronengine.SchedulerRequest{
        Ctx:         ctx,
        Spec:        "* * * * * *",
        WithLock:    true,
        AutoRelease: false,
        Usecase:     "runUsecase",
        Target:      "0",
        TTL:         5 * time.Second,
        Func: func() {
            // your job here
        },
    })
}
```

## Parameters
- Spec: Cron expression with seconds enabled (six fields). Example `* * * * * *` runs every second.
- WithLock: Enable distributed locking; if false, schedules without lock.
- AutoRelease: When true, releases the lock after job finishes; when false, rely on TTL expiration.
- Usecase: Logical name grouping locks (Preferably Function Name) (e.g., `send_email`).
- Target: Dimension to lock within the usecase (Which specific resource that job operates on) (e.g., user ID, date key).
- TTL: Lock expiration; prevents indefinite lock if a worker dies.
- Func: The function to execute on schedule.

Example Usecase Target: `send_email` (Usecase) and `2024-01-01` or `user:123` (Target)

## Lock-Only Example
Use the locking mechanism independently of cron to guard critical sections across workers.

```go
package main

import (
    "context"
    "log"
    "time"
    cronengine "github.com/bimonugraraga/radobaan-cron-engine"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()

    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatal(err)
    }

    lock := cronengine.NewLock(rdb)

    ok, err := lock.Acquire(cronengine.LockRequest{
        Ctx: ctx,
        Req: cronengine.Request{
            Usecase: "send_email",
            Target:  "user:123",
            TTL:     10 * time.Second,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    if !ok {
        log.Println("locked, skip")
        return
    }
    defer lock.Release(cronengine.LockRequest{
        Ctx: ctx,
        Req: cronengine.Request{
            Usecase: "send_email",
            Target:  "user:123",
        },
    })

    // protected work
}
```

## Code References
- Example: [main.go](file:///Users/bytedance/go/src/github.com/radobaan-cron-engine/examples/basic/main.go)
- API types: [request.go](file:///Users/bytedance/go/src/github.com/radobaan-cron-engine/request.go)
- Scheduler: [scheduler.go](file:///Users/bytedance/go/src/github.com/radobaan-cron-engine/scheduler.go)
- Locking: [lock.go](file:///Users/bytedance/go/src/github.com/radobaan-cron-engine/lock.go)
