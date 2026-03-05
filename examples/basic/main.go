package basic

import (
	"context"
	"fmt"
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
		Func:        runUsecase,
	})
}

func runUsecase() {
	fmt.Println("🔥 executing usecase...")

	fmt.Println("✅ done")
}
