package cronengine

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Client *redis.Client
}

func NewScheduler(rdb *redis.Client) *Scheduler {
	return &Scheduler{
		Client: rdb,
	}
}

func (s *Scheduler) Schedule(params SchedulerRequest) error {
	c := cron.New(
		cron.WithSeconds(),
	)

	if !params.WithLock {
		_, err := c.AddFunc(params.Spec, params.Func)
		return err
	}

	locker := NewLock(s.Client)

	_, err := c.AddFunc(params.Spec, func() {
		defer func() {
			if params.AutoRelease {

				locker.Release(LockRequest{
					Ctx: params.Ctx,
					Req: Request{
						Usecase: params.Usecase,
						Target:  params.Target,
					},
				})
			}
		}()

		ok, err := locker.Acquire(LockRequest{
			Ctx: params.Ctx,
			Req: Request{
				Usecase: params.Usecase,
				Target:  params.Target,
				TTL:     params.TTL,
			},
		})

		if err != nil {
			return
		}
		if !ok {
			return
		}
		params.Func()
	})
	c.Start()
	select {}
	return err
}
