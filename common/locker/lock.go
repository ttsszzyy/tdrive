/*
 * Author: 李鸿胤 leeyfann@gmail.com
 * Date: 2023-11-14 19:58:59
 * LastEditors: lihy lihy@zhiannet.com
 * Note: Need note condition
 */
package locker

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

var (
	defaultSpinTimeout = 30 * time.Second
	ErrLock            = errors.New("lock faild")
)

type ILocker interface {
	Lock() error
	Unlock() error
	LockCtx(ctx context.Context) error
	UnlockCtx(ctx context.Context) error
	SetExpire(seconds int)
	Spin(timeout time.Duration) error
	SpinCtx(ctx context.Context, timeout time.Duration) error
}

type RedisLocker struct {
	lk *redis.RedisLock
}

func NewRedisLocker(r *redis.Redis, key string) ILocker {
	return &RedisLocker{
		lk: redis.NewRedisLock(r, key),
	}
}

func (r *RedisLocker) Lock() error {
	f, e := r.lk.Acquire()
	if e != nil {
		return e
	}
	if !f {
		return ErrLock
	}
	return nil
}

func (r *RedisLocker) LockCtx(ctx context.Context) error {
	f, e := r.lk.AcquireCtx(ctx)
	if e != nil {
		return e
	}
	if !f {
		return ErrLock
	}
	return nil
}

func (r *RedisLocker) Unlock() error {
	_, e := r.lk.Release()
	return e
}

func (r *RedisLocker) UnlockCtx(ctx context.Context) error {
	_, e := r.lk.ReleaseCtx(ctx)
	return e
}

func (r *RedisLocker) SetExpire(seconds int) {
	r.lk.SetExpire(seconds)
}

// Spin 自旋同步等待
func (r *RedisLocker) Spin(timeout time.Duration) error {
	exp := time.Now().Add(timeout)
	tk := time.NewTicker(50 * time.Millisecond)
	defer tk.Stop()

	for {
		if time.Now().After(exp) {
			return ErrLock
		}
		if f, _ := r.lk.Acquire(); !f {
			<-tk.C
			continue
		}
		return nil
	}
}

// SpinCtx 含有ctx信号的自旋等待
func (r *RedisLocker) SpinCtx(ctx context.Context, timeout time.Duration) error {
	exp := time.Now().Add(timeout)
	tk := time.NewTicker(50 * time.Millisecond)
	defer tk.Stop()
loop:
	for {
		if time.Now().After(exp) {
			return ErrLock
		}
		if f, _ := r.lk.AcquireCtx(ctx); !f {
			select {
			case <-ctx.Done():
				return ErrLock
			case <-tk.C:
				continue loop
			}
		}
		return nil
	}
}
