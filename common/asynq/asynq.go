/*
 * Author: lihy lihy@zhiannet.com
 * Date: 2023-12-20 10:15:02
 * LastEditors: lihy lihy@zhiannet.com
 * Note: Need note condition
 */
package asynq

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// AsynqConf is the configuration struct for Asynq.
type AsynqConf struct {
	Addr         string `json:",default=127.0.0.1:6379"`
	Username     string `json:",optional"`
	Pass         string `json:",optional"`
	DB           int    `json:",optional,default=0"`
	Concurrency  int    `json:",optional,default=20"` // max concurrent process job task num
	SyncInterval int    `json:",optional,default=10"` // seconds, this field specifies how often sync should happen
}

// WithRedisConf sets redis configuration from RedisConf.
func (c *AsynqConf) WithRedisConf(r redis.RedisConf) *AsynqConf {
	c.Pass = r.Pass
	c.Addr = r.Host
	c.DB = 1
	return c
}

// WithOriginalRedisConf sets redis configuration from original RedisConf.
func (c *AsynqConf) WithOriginalRedisConf(r RedisConf) *AsynqConf {
	c.Pass = r.Pass
	c.Addr = r.Host
	c.Username = r.Username
	c.DB = r.Db
	return c
}

// NewRedisOpt returns a redis options from Asynq Configuration.
func (c *AsynqConf) NewRedisOpt() *asynq.RedisClientOpt {
	return &asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Pass,
		DB:       c.DB,
	}
}

// NewClient returns a client from the configuration.
func (c *AsynqConf) NewClient() *asynq.Client {
	return asynq.NewClient(c.NewRedisOpt())
}

// NewServer returns a worker from the configuration.
func (c *AsynqConf) NewServer(Queue map[string]int, concurrency int) *asynq.Server {
	return asynq.NewServer(
		c.NewRedisOpt(),
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("failed to exec asynq task, err : %+v \n", err)
				return true
			},
			Queues:      Queue,
			Concurrency: concurrency,
		},
	)
}

// NewScheduler returns a scheduler from the configuration.
func (c *AsynqConf) NewScheduler() *asynq.Scheduler {
	return asynq.NewScheduler(c.NewRedisOpt(), &asynq.SchedulerOpts{Location: time.Local})
}

// NewPeriodicTaskManager returns a periodic task manager from the configuration.
func (c *AsynqConf) NewPeriodicTaskManager(provider asynq.PeriodicTaskConfigProvider) *asynq.PeriodicTaskManager {
	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			SchedulerOpts:              &asynq.SchedulerOpts{Location: time.Local},
			RedisConnOpt:               c.NewRedisOpt(),
			PeriodicTaskConfigProvider: provider,                                    // this provider object is the interface to your config source
			SyncInterval:               time.Duration(c.SyncInterval) * time.Second, // this field specifies how often sync should happen
		})
	logx.Must(err)
	return mgr
}

// NewInspector returns a new instance of Asynq Inspector.
func (c *AsynqConf) NewInspector() *asynq.Inspector {
	return asynq.NewInspector(c.NewRedisOpt())
}
