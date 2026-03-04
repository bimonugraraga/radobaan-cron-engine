package cronengine

import (
	"context"
	"time"
)

// Lock Request
type Request struct {
	Usecase string        // Unique identifier for the usecase (example: "send_email")
	Target  string        // Target resource or entity to lock (example: "user:123", "2026-03-04")
	TTL     time.Duration // Time-to-live for the lock in seconds (default: 60)
}

type LockRequest struct {
	Ctx context.Context
	Req Request
}

// Cron Request
type SchedulerRequest struct {
	Ctx         context.Context
	Spec        string
	WithLock    bool
	AutoRelease bool
	Usecase     string
	Target      string
	TTL         time.Duration
	Func        func()
}
