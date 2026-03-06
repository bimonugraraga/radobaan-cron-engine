package cronengine

import (
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/redis/go-redis/v9"
)

type Lock struct {
	Client *redis.Client
}

func NewLock(rdb *redis.Client) *Lock {
	return &Lock{
		Client: rdb,
	}
}

func buildKey(usecase string) string {
	return fmt.Sprintf("lock:%s", usecase)
}

func hashTarget(target string) int64 {
	h := fnv.New64a()
	h.Write([]byte(target))
	return int64(h.Sum64()) % 1024
}

func (l *Lock) Acquire(params LockRequest) (bool, error) {
	if params.Req.Usecase == "" || params.Req.Target == "" {
		return false, errors.New("usecase and target are required")
	}

	key := buildKey(params.Req.Usecase)
	offset := hashTarget(params.Req.Target)
	if params.Req.TTL <= 0 {
		params.Req.TTL = 60
	}
	ttl := int64(params.Req.TTL.Seconds())

	res, err := AcquireScript.Run(
		params.Ctx,
		l.Client,
		[]string{key},
		offset,
		ttl).Int()
	if err != nil {
		return false, err
	}
	if res == 0 {
		return false, nil
	}

	return true, nil
}

func (l *Lock) Release(params LockRequest) error {
	if params.Req.Usecase == "" || params.Req.Target == "" {
		return errors.New("usecase and target are required")
	}

	return l.Client.SetBit(params.Ctx, buildKey(params.Req.Usecase), hashTarget(params.Req.Target), 0).Err()
}
